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
	client    *mongo.Client
	log       *zap.Logger
	cfg       *config.MongoDBConnectionConfig
	connected bool
}

func NewMongoDB(logger *zap.Logger, cfg *config.MongoDBConnectionConfig) *MongoDB {
	return &MongoDB{
		client:    nil,
		log:       logger,
		cfg:       cfg,
		connected: false,
	}
}

func (m *MongoDB) Connect() error {
	var err error
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", m.cfg.Username, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.Database)

	opts := options.Client().ApplyURI(uri)
	m.client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	var result bson.M
	if err := m.client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		return err
	}

	m.connected = true
	return nil
}

func (m *MongoDB) IsConnected() bool {
	return m.connected
}

func (m *MongoDB) Release() {
	if m.client != nil {
		m.client.Disconnect(context.TODO())
	}
}

func (m *MongoDB) GetDB() *mongo.Database {
	return m.client.Database(m.cfg.Database)
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.GetDB().Collection(name)
}
