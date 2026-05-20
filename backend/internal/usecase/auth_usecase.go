package usecase

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repositories"
	"backend/utils"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUserUsecase interface {
	Register(email, password, name, phone string, role entities.UserRole) (*entities.UserResponse, error)
	Login(email, password string) (*entities.UserResponse, string, error)
	GetProfile(userID uuid.UUID) (*entities.UserResponse, error)
}

type authUserUseCase struct {
	userRepo repositories.UserRepository
	jwtUtil  *utils.JWTUtil
}

func NewAuthUserUsecase(
	userRepo repositories.UserRepository,
	jwtUtil *utils.JWTUtil,
) AuthUserUsecase {
	return &authUserUseCase{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (u *authUserUseCase) Register(email, password, name, phone string, role entities.UserRole) (*entities.UserResponse, error) {

	_, err := u.userRepo.GetByEmail(email)
	if err == nil {
		return nil, errors.New("user with this email already exists in this tenant")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate username from email
	username := email[:strings.Index(email, "@")]

	// Create user
	user := &entities.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
		FullName: name,
		Phone:    phone,
		Role:     role,
		IsActive: true,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user.ToResponse(), nil
}

func (u *authUserUseCase) Login(email, password string) (*entities.UserResponse, string, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("email atau password salah")
		}
		return nil, "", err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", errors.New("account is inactive")
	}

	// Verify password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("email atau password salah")
	}

	token, err := u.jwtUtil.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", errors.New("ulangi beberapa saat lagi")
	}

	return user.ToResponse(), token, nil
}

func (u *authUserUseCase) GetProfile(userID uuid.UUID) (*entities.UserResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user.ToResponse(), nil
}
