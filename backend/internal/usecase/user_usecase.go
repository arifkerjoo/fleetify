package usecase

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repositories"
	"backend/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUsecase interface {
	GetAllUsers(limit, offset int, search string) ([]entities.UserListResponse, int64, error)
	GetUserByID(userID uuid.UUID) (*entities.UserResponse, error)
	UpdateUser(userID uuid.UUID, req UpdateUserRequest) (*entities.UserResponse, error)
	UpdateProfileImage(userID uuid.UUID, imagePath string) error
	DeleteUser(userID, requestUserID uuid.UUID) error
	ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
}

type userUsecase struct {
	userRepo repositories.UserRepository
}

func NewUserUsecase(userRepo repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetAllUsers(limit, offset int, search string) ([]entities.UserListResponse, int64, error) {
	users, total, err := u.userRepo.GetAllUsers(limit, offset, search)
	if err != nil {
		return nil, 0, errors.New("failed to fetch users")
	}

	responses := make([]entities.UserListResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToListResponse()
	}

	return responses, total, nil
}

func (u *userUsecase) GetUserByID(userID uuid.UUID) (*entities.UserResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user.ToResponse(), nil
}

func (u *userUsecase) UpdateUser(userID uuid.UUID, req UpdateUserRequest) (*entities.UserResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check if username already taken (if changed)
	if req.Username != "" && req.Username != user.Username {
		existing, _ := u.userRepo.GetByUsername(req.Username)
		if existing != nil && existing.ID != userID {
			return nil, errors.New("username already taken")
		}
		user.Username = req.Username
	}

	// Update fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	return user.ToResponse(), nil
}

func (u *userUsecase) UpdateProfileImage(userID uuid.UUID, imagePath string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := u.userRepo.Update(user); err != nil {
		return errors.New("failed to update profile image")
	}

	return nil
}

func (u *userUsecase) DeleteUser(userID, requestUserID uuid.UUID) error {
	// Prevent self-deletion
	if userID == requestUserID {
		return errors.New("you cannot delete your own account")
	}

	_, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := u.userRepo.Delete(userID); err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}

func (u *userUsecase) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Verify old password
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.Password = hashedPassword

	if err := u.userRepo.Update(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
