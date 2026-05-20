package repositories

import (
	"backend/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entities.User) error
	GetByEmail(email string) (*entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	GetByID(id uuid.UUID) (*entities.User, error)
	GetAllUsers(limit, offset int, search string) ([]entities.User, int64, error)
	Update(user *entities.User) error
	Delete(id uuid.UUID) error
}
