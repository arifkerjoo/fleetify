package repositories

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type masterItemRepository struct {
	db *gorm.DB
}

func NewMasterItemRepository(db *gorm.DB) repositories.MasterItemRepository {
	return &masterItemRepository{db: db}
}

func (r *masterItemRepository) Create(masterItem *entities.MasterItem) error {
	return r.db.Create(masterItem).Error
}

func (r *masterItemRepository) GetByID(id uuid.UUID) (*entities.MasterItem, error) {
	var masterItem entities.MasterItem
	if err := r.db.Where("id = ?", id).First(&masterItem).Error; err != nil {
		return nil, err
	}
	return &masterItem, nil
}

func (r *masterItemRepository) Update(masterItem *entities.MasterItem) error {
	return r.db.Save(masterItem).Error
}

func (r *masterItemRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.MasterItem{}, id).Error
}

func (r *masterItemRepository) GetAllItems(limit, offset int, search string) ([]entities.MasterItem, int64, error) {
	var masterItems []entities.MasterItem
	var total int64

	query := r.db.Model(&entities.MasterItem{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"item_name ILIKE ? OR item_type ILIKE ?",
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
		Find(&masterItems).Error; err != nil {
		return nil, 0, err
	}

	return masterItems, total, nil
}
