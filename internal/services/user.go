package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sonochiwa/auth/internal/mongodb"
	"github.com/sonochiwa/auth/internal/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	log         *zap.Logger
	mongoClient *mongodb.MongoDB
	dbClient    *sqlx.DB
	collection  string
}

func NewUserService(logger *zap.Logger, mongoClient *mongodb.MongoDB, db *sqlx.DB) *UserService {
	return &UserService{
		log:         logger,
		mongoClient: mongoClient,
		dbClient:    db,
		collection:  "users",
	}
}

var UserAlreadyExistsError = fmt.Errorf("user already exists")

func (s *UserService) GetByGuid(guid string) (*models.User, error) {
	var user *models.User
	collection := s.mongoClient.GetCollection("users")
	err := collection.FindOne(context.TODO(), bson.M{"guid": guid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByLogin(login string) (*models.User, error) {
	var user *models.User
	collection := s.mongoClient.GetCollection("users")
	err := collection.FindOne(context.TODO(), bson.M{"login": login}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Register(login string, password string) (*models.User, error) {
	var user *models.User
	collections := s.mongoClient.GetCollection(s.collection)
	existedUser, err := s.GetByLogin(login)
	if existedUser != nil {
		return nil, UserAlreadyExistsError
	}

	if err != nil {
		return nil, err
	}

	userGuid := uuid.New().String()
	newUser := &models.User{
		Login:     login,
		GUID:      userGuid,
		LoginType: models.LoginType{ID: 1, Name: "email"},
		Name:      "",
		LastName:  "",
		CreatedAt: time.Now(),
	}

	user = &models.User{
		Login:     login,
		CreatedAt: time.Now(),
	}

	insertResult, err := collections.InsertOne(context.TODO(), newUser)
	if err != nil {
		return nil, err
	}

	sql := `INSERT INTO passwords (user_guid, password, expired_at) VALUES ($1, $2, $3)`
	hashedPwd, err := s.HashPassword(password)

	if err != nil {
		s.log.Error("Error while hashing password", zap.Error(err))
		err2 := s.DeleteById(insertResult.InsertedID.(primitive.ObjectID))
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	_, err = s.dbClient.Exec(sql, userGuid, hashedPwd, time.Now().Add(time.Hour*24))
	if err != nil {
		s.log.Error("Error while saving password", zap.Error(err))
		err2 := s.DeleteById(insertResult.InsertedID.(primitive.ObjectID))
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteById(id primitive.ObjectID) error {
	collection := s.mongoClient.GetCollection(s.collection)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *UserService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
