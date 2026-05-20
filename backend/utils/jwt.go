package utils

import (
	"backend/internal/domain/entities"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID         `json:"user_id"`
	Email  string            `json:"email"`
	Role   entities.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secretKey []byte
}

func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{
		secretKey: []byte(secretKey),
	}
}

func (j *JWTUtil) GenerateToken(userID uuid.UUID, email string, role entities.UserRole) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWTUtil) RefreshToken(oldToken string) (string, error) {
	claims, err := j.ValidateToken(oldToken)
	if err != nil {
		return "", err
	}

	return j.GenerateToken(claims.UserID, claims.Email, claims.Role)
}
