package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sonochiwa/auth/internal/config"
	"github.com/sonochiwa/auth/internal/db"
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
	for {
		if mongoClient.IsConnected() {
			logger.Info("Connecting to MongoDB")
			break
		}
		time.Sleep(1 * time.Second)
	}

	logger.Info("Connecting to Postgres")
	dbClient, err := db.NewDBConnection(&cfg.DB)
	if err != nil {
		logger.Fatal("Failed to connect to Postgres", zap.Error(err))
	}

	row := dbClient.QueryRow("select now() as t")
	if row.Err() != nil {
		logger.Error("Failed to get current time", zap.Error(err))
	}
	var t time.Time
	err = row.Scan(&t)
	if err != nil {
		logger.Error("Failed to get current time", zap.Error(err))
	}
	logger.Info("Current DB Time", zap.String("Time", t.Format("02.01.2006 15:04:05")))

	err = db.ApplyMigrations(cfg.DB.Type, dbClient)
	if err != nil {
		logger.Fatal("Failed to apply migrations", zap.Error(err))
	}

	logger.Info("Migrations applied")

	wApp := web.NewWebServer(logger, &cfg.Web)
	registerRoutes(logger, mongoClient, dbClient, wApp)
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

func registerRoutes(logger *zap.Logger, mongoClient *mongodb.MongoDB, db *sqlx.DB, wApp *web.WebServer) {
	logger.Info("Register routes")

	userService := services.NewUserService(logger, mongoClient, db)

	wApp.RegisterRoutes([]controllers.Controller{
		controllers.NewAuthController(logger, userService),
		controllers.NewRegisterController(logger, userService),
	})
}
