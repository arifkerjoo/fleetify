package routes

import (
	"backend/internal/interfaces/handlers"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

type Handlers struct {
	Auth *handlers.AuthHandler
	User *handlers.UserHandler
}

func SetupRoutes(app *fiber.App, h *Handlers, jwtUtil *utils.JWTUtil) {

	api := app.Group("/api/v1")

	AuthRoutes(api, h.Auth, jwtUtil)
	h.User.RegisterRoutes(api)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
			"version": "1.0.0",
		})
	})

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    404,
				"message": "Route not found",
			},
			"success": false,
		})
	})
}
