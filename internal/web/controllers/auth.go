package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sonochiwa/auth/internal/mongodb/models"
	"github.com/sonochiwa/auth/internal/services"
	"go.uber.org/zap"
)

var LoginRequiredError = fmt.Errorf("login required")
var PasswordRequiredError = fmt.Errorf("password required")

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r AuthRequest) Validate() (bool, error) {
	if len(r.Login) == 0 {
		return false, LoginRequiredError
	}

	if len(r.Password) == 0 {
		return false, PasswordRequiredError
	}

	return true, nil
}

type AuthResponse struct {
	Ok           bool         `json:"ok"`
	Token        string       `json:"token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	User         *models.User `json:"user,omitempty"`
	Cause        string       `json:"cause,omitempty"`
}

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

func (c *AuthController) GetGroup() string {
	return "/auth"
}

func (c *AuthController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method:  "GET",
			Path:    "/auth",
			Handler: c.authHandler(),
		},
	}
}

func (c *AuthController) authHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var req AuthRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		if valid, err := AuthRequest.Validate(req); !valid {
			return c.JSON(AuthResponse{
				Ok:    false,
				Cause: err.Error(),
			})
		}

		return c.JSON(AuthResponse{
			Ok:           true,
			Token:        "123",
			RefreshToken: "123",
			User: &models.User{
				Login: req.Login,
			},
		})
	}
}
