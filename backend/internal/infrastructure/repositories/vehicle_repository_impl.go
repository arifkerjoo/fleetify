package repositories

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) repositories.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(vehicle *entities.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *vehicleRepository) GetByID(id uuid.UUID) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle
	if err := r.db.Where("id = ?", id).First(&vehicle).Error; err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *vehicleRepository) Update(vehicle *entities.Vehicle) error {
	return r.db.Save(vehicle).Error
}

func (r *vehicleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Vehicle{}, id).Error
}

func (r *vehicleRepository) GetAllVehicles(limit, offset int, search string) ([]entities.Vehicle, int64, error) {
	var vehicles []entities.Vehicle
	var total int64

	query := r.db.Model(&entities.Vehicle{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"brand ILIKE ? OR model ILIKE ? OR license_plate ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&vehicles).Error; err != nil {
		return nil, 0, err
	}

	return vehicles, total, nil
}
