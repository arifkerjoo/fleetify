package handlers

import (
	"backend/internal/domain/entities"
	"backend/internal/interfaces/middleware"
	"backend/internal/usecase"
	"backend/utils"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUserUseCase usecase.AuthUserUsecase
	jwtUtil         *utils.JWTUtil
}

type RegisterUserRequest struct {
	Name     string            `json:"name" validate:"required,min=3,max=100"`
	Email    string            `json:"email" validate:"required,email"`
	Phone    string            `json:"phone" validate:"required,min=10,max=20"`
	Password string            `json:"password" validate:"required,min=8"`
	Role     entities.UserRole `json:"role" validate:"omitempty,oneof=admin manager user"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Message string            `json:"message"`
	User    interface{}       `json:"user,omitempty"`
	Token   string            `json:"token,omitempty"`
	Role    entities.UserRole `json:"role,omitempty"`
}

func NewAuthHandler(authUserUseCase usecase.AuthUserUsecase, jwtUtil *utils.JWTUtil) *AuthHandler {
	return &AuthHandler{
		authUserUseCase: authUserUseCase,
		jwtUtil:         jwtUtil,
	}
}

func (h *AuthHandler) RegisterRoutes(api fiber.Router) {
	auth := api.Group("/user")

	// Public routes
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)

	// Protected routes
	protected := auth.Group("")
	protected.Use(middleware.JWTMiddleware(h.jwtUtil))
	{
		protected.Get("/profile", h.GetProfile)
		protected.Post("/refresh-token", h.RefreshToken)
		protected.Post("/logout", h.Logout)
		protected.Get("/session", h.GetSession) // For Next.js to check session
	}

	// Admin only routes
	admin := auth.Group("/admin")
	admin.Use(middleware.JWTMiddleware(h.jwtUtil))
	admin.Use(middleware.RequireRoleMiddleware("admin"))
	{
		admin.Get("/users", h.ListUsers)
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, token, err := h.authUserUseCase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Determine if running in production (HTTPS)
	isProduction := os.Getenv("ENV") == "production"

	// Set HTTPOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token_base",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   isProduction, // true in production (HTTPS)
		SameSite: "Lax",        // "None" if cross-domain, "Lax" for same-domain
		MaxAge:   86400,        // 24 hours
	})

	return c.JSON(AuthResponse{
		Message: "Login successful",
		User:    user,
		Token:   token,
		Role:    user.Role,
	})
}

func (h *AuthHandler) GetSession(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	role := c.Locals("role").(string)

	user, err := h.authUserUseCase.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"authenticated": true,
		"user": fiber.Map{
			"id":    userID,
			"email": user.Email,
			"name":  user.FullName,
			"role":  role,
		},
	})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	// Clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token_base",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		MaxAge:   -1, // Delete cookie
	})

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) GetProfile(c fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.authUserUseCase.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile retrieved successfully",
		"user":    user,
	})
}

func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	token := c.Cookies("auth_token_base")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && len(authHeader) > 7 {
			token = authHeader[7:]
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token not found",
		})
	}

	newToken, err := h.jwtUtil.RefreshToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to refresh token",
		})
	}

	// Update cookie
	isProduction := os.Getenv("ENV") == "production"
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token_base",
		Value:    newToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   isProduction,
		SameSite: "Lax",
		MaxAge:   86400,
	})

	return c.JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"token":   newToken,
	})
}

func (h *AuthHandler) ListUsers(c fiber.Ctx) error {
	// Only accessible by admin (protected by middleware)
	return c.JSON(fiber.Map{
		"message": "Admin only endpoint",
	})
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user, err := h.authUserUseCase.Register(req.Email, req.Password, req.Name, req.Phone, role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Message: "User registered successfully",
		User:    user,
	})
}
