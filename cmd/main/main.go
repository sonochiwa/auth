package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sonochiwa/auth/internal/config"
	"github.com/sonochiwa/auth/internal/logging"
	"github.com/sonochiwa/auth/internal/mongodb"
	"github.com/sonochiwa/auth/internal/services"
	"github.com/sonochiwa/auth/internal/web"
	"github.com/sonochiwa/auth/internal/web/controllers"
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
	logger.Info("Connecting to MongoDB")
	wApp := web.NewWebServer(logger, &cfg.Web)
	registerRoutes(logger, mongoClient, wApp)
	go wApp.Run()

	// Graceful shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	takeSig := <-sigChan
	logger.Info("Shutting down graceful", zap.String("signal", takeSig.String()))
}
func runWebServer(logger *zap.Logger, cfg *config.WebServerConfig) {
	logger.Info("Starting web server")
	app := fiber.New()
	app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func registerRoutes(logger *zap.Logger, mongoClient *mongodb.MongoDB, wApp *web.WebServer) {
	logger.Info("Register routes")

	userService := services.NewUserService(logger, mongoClient)

	wApp.RegisterRoutes([]controllers.Controller{
		controllers.NewAuthController(logger, userService),
	})
}
