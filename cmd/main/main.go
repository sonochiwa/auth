package main

import (
	"github.com/sonochiwa/auth/internal/config"
	"github.com/sonochiwa/auth/internal/logging"
	"github.com/sonochiwa/auth/internal/mongodb"
	"go.uber.org/zap"
	"os"
	"os/signal"

	"math/rand"
	"time"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func main() {

	cfg := config.InitConfiguration()
	logger := logging.NewLogger("auth.log")
	mongoClient := mongodb.NewMongoDB(logger, &cfg.Mongo)
	defer mongoClient.Release()
	err := mongoClient.Connect()
	if err != nil {
		logger.Fatal("Failed connect to MongoDB!", zap.Error(err))
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	takeSig := <-sigChan
	logger.Info("Shutting down graceful", zap.String("signal", takeSig.String()))
}
