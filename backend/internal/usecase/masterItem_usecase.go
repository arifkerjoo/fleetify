package usecase

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repositories"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MasterItemUsecase interface {
	CreateMasterItem(req CreateMasterItemRequest) (*entities.MasterItemResponse, error)
	GetMasterItemByID(id uuid.UUID) (*entities.MasterItemResponse, error)
	GetAllMasterItems(limit, offset int, search string) ([]entities.MasterItemResponse, int64, error)
	UpdateMasterItem(id uuid.UUID, req UpdateMasterItemRequest) (*entities.MasterItemResponse, error)
	DeleteMasterItem(id uuid.UUID) error
}

type CreateMasterItemRequest struct {
	ItemName string            `json:"item_name" binding:"required"`
	ItemType entities.ItemType `json:"item_type" binding:"required"`
	Price    float64           `json:"price" binding:"required"`
	IsActive bool              `json:"is_active"`
}

type UpdateMasterItemRequest struct {
	ItemName string  `json:"item_name"`
	Price    float64 `json:"price"`
	IsActive *bool   `json:"is_active"`
}

type masterItemUsecase struct {
	masterItemRepo repositories.MasterItemRepository
}

func NewMasterItemUsecase(
	masterItemRepo repositories.MasterItemRepository,
) MasterItemUsecase {
	return &masterItemUsecase{
		masterItemRepo: masterItemRepo,
	}
}

func (u *masterItemUsecase) CreateMasterItem(
	req CreateMasterItemRequest,
) (*entities.MasterItemResponse, error) {

	masterItem := &entities.MasterItem{
		BaseModel: entities.BaseModel{
			ID: uuid.New(),
		},
		ItemName: req.ItemName,
		ItemType: req.ItemType,
		Price:    req.Price,
		IsActive: req.IsActive,
	}

	if err := u.masterItemRepo.Create(masterItem); err != nil {
		return nil, errors.New("failed to create master item")
	}

	return masterItem.ToResponse(), nil
}

func (u *masterItemUsecase) GetMasterItemByID(
	id uuid.UUID,
) (*entities.MasterItemResponse, error) {

	masterItem, err := u.masterItemRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("master item not found")
		}

		return nil, err
	}

	return masterItem.ToResponse(), nil
}

func (u *masterItemUsecase) GetAllMasterItems(
	limit, offset int,
	search string,
) ([]entities.MasterItemResponse, int64, error) {

	masterItems, total, err := u.masterItemRepo.GetAllItems(
		limit,
		offset,
		search,
	)

	if err != nil {
		return nil, 0, errors.New("failed to fetch master items")
	}

	responses := make([]entities.MasterItemResponse, len(masterItems))

	for i, masterItem := range masterItems {
		responses[i] = *masterItem.ToResponse()
	}

	return responses, total, nil
}

func (u *masterItemUsecase) UpdateMasterItem(
	id uuid.UUID,
	req UpdateMasterItemRequest,
) (*entities.MasterItemResponse, error) {

	masterItem, err := u.masterItemRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("master item not found")
		}

		return nil, err
	}

	if req.ItemName != "" {
		masterItem.ItemName = req.ItemName
	}

	if req.Price > 0 {
		masterItem.Price = req.Price
	}

	if req.IsActive != nil {
		masterItem.IsActive = *req.IsActive
	}

	if err := u.masterItemRepo.Update(masterItem); err != nil {
		return nil, errors.New("failed to update master item")
	}

	return masterItem.ToResponse(), nil
}

func (u *masterItemUsecase) DeleteMasterItem(id uuid.UUID) error {
	_, err := u.masterItemRepo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("master item not found")
		}

		return err
	}

	if err := u.masterItemRepo.Delete(id); err != nil {
		return errors.New("failed to delete master item")
	}

	return nil
}
