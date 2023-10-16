package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sonochiwa/auth/internal/config"
	"github.com/sonochiwa/auth/internal/web/controllers"
	"go.uber.org/zap"
	"reflect"
)

type WebServer struct {
	log    *zap.Logger
	cfg    *config.WebServerConfig
	client *fiber.App
}

func NewWebServer(logger *zap.Logger, cfg *config.WebServerConfig) *WebServer {
	return &WebServer{
		log:    logger,
		cfg:    cfg,
		client: fiber.New(fiber.Config{}),
	}
}

func (w *WebServer) RegisterRoutes(routes []controllers.Controller) {
	for _, route := range routes {
		w.log.Error(
			"Register controller",
			zap.String("controller", reflect.TypeOf(route).Elem().Name()),
			zap.String("path", route.GetPath()),
			zap.String("method", route.GetMethod()),
		)
		switch route.GetMethod() {
		case "GET":
			w.client.Get(route.GetPath(), route.GetHandler())
		case "POST":
			w.client.Post(route.GetPath(), route.GetHandler())
		default:
			w.log.Error(
				"Unsupported method",
				zap.String("controller", reflect.TypeOf(route).Name()),
				zap.String("path", route.GetPath()),
				zap.String("method", route.GetMethod()),
			)
		}
	}
}

func (w *WebServer) Run() {
	w.log.Info("Starting web server")
	w.client.Listen(fmt.Sprintf("0.0.0.0:%d", w.cfg.Port))
}
