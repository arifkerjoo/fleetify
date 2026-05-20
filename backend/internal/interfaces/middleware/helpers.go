package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func GetUserIDFromContext(c fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found")
	}
	return userID, nil
}

func GetRoleFromContext(c fiber.Ctx) string {
	role, ok := c.Locals("role").(string)
	if !ok {
		return ""
	}
	return role
}
