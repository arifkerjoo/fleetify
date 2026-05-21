package repositories

import (
	"backend/internal/domain/entities"

	"github.com/google/uuid"
)

type MasterItemRepository interface {
	Create(masterItem *entities.MasterItem) error
	GetByID(id uuid.UUID) (*entities.MasterItem, error)
	GetAllItems(limit, offset int, search string) ([]entities.MasterItem, int64, error)
	Update(masterItem *entities.MasterItem) error
	Delete(id uuid.UUID) error
}
