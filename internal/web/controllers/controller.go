package controllers

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetMethod() string // return http method
	GetPath() string   // return path
	GetHandler() func(c *fiber.Ctx) error
}
