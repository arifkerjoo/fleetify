package utils

import (
	"errors"
	"math"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// Pembulatan
func CeilToN(x float64, n int) float64 {
	factor := math.Pow(10, float64(n))
	return math.Ceil(x*factor) / factor
}

func RemoveBaseURL(fullURL string) string {
	u, _ := url.Parse(fullURL)
	clean := strings.TrimPrefix(u.Path, "/uploads")
	return clean
}

func GetUserIDFromContext(c fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return userID, nil
}

func GetEmailFromContext(c fiber.Ctx) (string, error) {
	email, ok := c.Locals("email").(string)
	if !ok {
		return "", errors.New("email not found in context")
	}
	return email, nil
}

func GetAuthContext(c fiber.Ctx) (userID uuid.UUID, email string, err error) {
	userID, err = GetUserIDFromContext(c)
	if err != nil {
		return
	}

	email, err = GetEmailFromContext(c)
	return
}
