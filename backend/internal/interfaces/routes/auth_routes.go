package routes

import (
	"backend/internal/interfaces/handlers"
	"backend/internal/interfaces/middleware"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

func AuthRoutes(r fiber.Router, authHandler *handlers.AuthHandler, jwtUtil *utils.JWTUtil) {
	auth := r.Group("/auth")

	// Public routes
	setupPublicAuthRoutes(auth, authHandler)

	// Protected routes
	setupProtectedAuthRoutes(auth, authHandler, jwtUtil)

	// Admin routes
	setupAdminAuthRoutes(auth, authHandler, jwtUtil)
}

// setupPublicAuthRoutes sets up public authentication routes
func setupPublicAuthRoutes(auth fiber.Router, handler *handlers.AuthHandler) {
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
}

// setupProtectedAuthRoutes sets up protected authentication routes
func setupProtectedAuthRoutes(auth fiber.Router, handler *handlers.AuthHandler, jwtUtil *utils.JWTUtil) {
	protected := auth.Group("")
	protected.Use(middleware.JWTMiddleware(jwtUtil))

	// Profile management
	protected.Get("/me", handler.GetProfile)
	protected.Get("/profile", handler.GetProfile)

	// Session management
	protected.Get("/session", handler.GetSession)
	protected.Post("/refresh-token", handler.RefreshToken)
	protected.Post("/logout", handler.Logout)
}

// setupAdminAuthRoutes sets up admin-only authentication routes
func setupAdminAuthRoutes(auth fiber.Router, handler *handlers.AuthHandler, jwtUtil *utils.JWTUtil) {
	admin := auth.Group("/admin")
	admin.Use(middleware.JWTMiddleware(jwtUtil))
	admin.Use(middleware.RequireRoleMiddleware("admin"))

	admin.Get("/users", handler.ListUsers)
}
