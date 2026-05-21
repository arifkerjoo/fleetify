package repositories

import (
	"backend/internal/domain/entities"

	"github.com/google/uuid"
)

type VehicleRepository interface {
	Create(vehicle *entities.Vehicle) error
	GetByID(id uuid.UUID) (*entities.Vehicle, error)
	GetAllVehicles(limit, offset int, search string) ([]entities.Vehicle, int64, error)
	Update(vehicle *entities.Vehicle) error
	Delete(id uuid.UUID) error
}
