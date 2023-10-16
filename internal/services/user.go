package services

import (
	"context"
	"github.com/sonochiwa/auth/internal/mongodb"
	"github.com/sonochiwa/auth/internal/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type UserService struct {
	log         *zap.Logger
	mongoClient *mongodb.MongoDB
	collection  string
}

func NewUserService(logger *zap.Logger, mongoClient *mongodb.MongoDB) *UserService {
	return &UserService{
		log:         logger,
		mongoClient: mongoClient,
		collection:  "users",
	}
}

func (s *UserService) GetByGuid(guid string) (*models.User, error) {
	var user *models.User
	collection := s.mongoClient.GetCollection("users")
	err := collection.FindOne(context.TODO(), bson.M{"guid": guid}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByLogin(login string) (*models.User, error) {
	var user *models.User
	collection := s.mongoClient.GetCollection("users")
	err := collection.FindOne(context.TODO(), bson.M{"login": login}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
