package middleware

import (
	"backend/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// JWTMiddleware validates JWT token and stores claims in context
func JWTMiddleware(jwtUtil *utils.JWTUtil) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Try to get token from cookie first
		tokenString := c.Cookies("auth_token_base")

		// If not in cookie, check Authorization header
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Authentication required",
				})
			}

			tokenString = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
			if tokenString == "" || tokenString == authHeader {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid authorization header format",
				})
			}
		}

		claims, err := jwtUtil.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		if claims.UserID == uuid.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}

		// Store claims in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RequireRoleMiddleware checks if user has required role
func RequireRoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Role not found in token",
			})
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}
