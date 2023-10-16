package mongodb

import (
	"context"
	"fmt"
	"github.com/sonochiwa/auth/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.uber.org/zap"
)

type MongoDB struct {
	client *mongo.Client
	log    *zap.Logger
	cfg    *config.MongoDBConnectionConfig
}

func NewMongoDB(logger *zap.Logger, cfg *config.MongoDBConnectionConfig) *MongoDB {
	return &MongoDB{
		client: nil,
		log:    logger,
		cfg:    cfg,
	}
}

func (m *MongoDB) Connect() error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", m.cfg.Username, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.Database)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) Release() {
	if m.client != nil {
		m.client.Disconnect(context.TODO())
	}
}
