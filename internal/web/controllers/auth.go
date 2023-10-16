package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sonochiwa/auth/internal/services"
	"go.uber.org/zap"
)

type AuthController struct {
	log         *zap.Logger
	userService *services.UserService
}

func NewAuthController(log *zap.Logger, userService *services.UserService) *AuthController {
	return &AuthController{
		log:         log,
		userService: userService,
	}
}

func (c *AuthController) GetMethod() string {
	return "GET"
}

func (c *AuthController) GetPath() string {
	return "/auth"
}

func (c *AuthController) GetHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("Auth")
	}
}
