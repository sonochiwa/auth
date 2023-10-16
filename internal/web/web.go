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
		group := w.client.Group(route.GetGroup())
		for _, handler := range route.GetHandlers() {
			switch handler.GetMethod() {
			case "GET":
				group.Get(handler.GetPath(), handler.GetHandler())
			case "POST":
				group.Post(handler.GetPath(), handler.GetHandler())
			case "PUT":
				group.Put(handler.GetPath(), handler.GetHandler())
			case "PATCH":
				group.Patch(handler.GetPath(), handler.GetHandler())
			case "DELETE":
				group.Delete(handler.GetPath(), handler.GetHandler())

			default:
				w.log.Error("Unsupported method",
					zap.String("controller", reflect.TypeOf(route).Elem().Name()),
					zap.String("path", handler.GetPath()),
					zap.String("method", handler.GetMethod()),
				)
			}
		}
	}
}

func (w *WebServer) Run() {
	w.log.Info("Starting web server")
	w.client.Listen(fmt.Sprintf("0.0.0.0:%d", w.cfg.Port))
}
