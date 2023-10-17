package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sonochiwa/auth/internal/mongodb/models"
	"github.com/sonochiwa/auth/internal/services"
	"go.uber.org/zap"
)

var PasswordNotEqualError = fmt.Errorf("password not equal")

type RegisterRequest struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (r RegisterRequest) Validate() (bool, error) {
	if len(r.Login) == 0 {
		return false, LoginRequiredError
	}

	if len(r.Password) == 0 {
		return false, PasswordRequiredError
	}

	if len(r.ConfirmPassword) == 0 {
		return false, PasswordRequiredError
	}

	if r.Password != r.ConfirmPassword {
		return false, PasswordRequiredError
	}

	return true, nil
}

type RegisterResponse struct {
	Ok    bool         `json:"ok"`
	User  *models.User `json:"user,omitempty"`
	Cause string       `json:"cause,omitempty"`
}

type RegisterController struct {
	log         *zap.Logger
	userService *services.UserService
}

func NewRegisterController(log *zap.Logger, userService *services.UserService) *RegisterController {
	return &RegisterController{
		log:         log,
		userService: userService,
	}
}

func (c *RegisterController) GetGroup() string {
	return "/auth"
}

func (c *RegisterController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method:  "POST",
			Path:    "/register",
			Handler: c.registerHandler(),
		},
	}
}

func (c *RegisterController) registerHandler() func(*fiber.Ctx) error {
	return func(fc *fiber.Ctx) error {
		fc.Accepts("application/json")
		var req RegisterRequest
		if err := fc.BodyParser(&req); err != nil {
			return err
		}

		if valid, err := RegisterRequest.Validate(req); !valid {
			return fc.JSON(RegisterResponse{
				Ok:    false,
				Cause: err.Error(),
			})
		}

		user, err := c.userService.Register(req.Login, req.Password)
		if err != nil {
			return fc.JSON(RegisterResponse{
				Ok:    false,
				Cause: err.Error(),
			})
		}

		return fc.JSON(RegisterResponse{
			Ok:   true,
			User: user,
		})
	}
}
