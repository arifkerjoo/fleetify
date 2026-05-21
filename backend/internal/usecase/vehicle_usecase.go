package usecase

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repositories"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleUsecase interface {
	CreateVehicle(req CreateVehicleRequest) (*entities.VehicleResponse, error)
	GetVehicleByID(id uuid.UUID) (*entities.VehicleResponse, error)
	GetAllVehicles(limit, offset int, search string) ([]entities.VehicleResponse, int64, error)
	UpdateVehicle(id uuid.UUID, req UpdateVehicleRequest) (*entities.VehicleResponse, error)
	DeleteVehicle(id uuid.UUID) error
}

type CreateVehicleRequest struct {
	LicensePlate string                 `json:"license_plate" binding:"required"`
	Model        string                 `json:"model" binding:"required"`
	Brand        string                 `json:"brand" binding:"required"`
	Year         int                    `json:"year" binding:"required"`
	Status       entities.VehicleStatus `json:"status" binding:"required"`
}

type UpdateVehicleRequest struct {
	LicensePlate string `json:"license_plate"`
	Model        string `json:"model"`
}

type vehicleUsecase struct {
	vehicleRepo repositories.VehicleRepository
}

func NewVehicleUsecase(vehicleRepo repositories.VehicleRepository) VehicleUsecase {
	return &vehicleUsecase{
		vehicleRepo: vehicleRepo,
	}
}

func (u *vehicleUsecase) CreateVehicle(req CreateVehicleRequest) (*entities.VehicleResponse, error) {
	vehicle := &entities.Vehicle{
		BaseModel: entities.BaseModel{
			ID: uuid.New(),
		},
		LicensePlate: req.LicensePlate,
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         req.Year,
		Status:       req.Status,
	}

	if err := u.vehicleRepo.Create(vehicle); err != nil {
		return nil, errors.New("Create vehicle data failed")
	}

	return vehicle.ToResponse(), nil
}

func (u *vehicleUsecase) GetVehicleByID(id uuid.UUID) (*entities.VehicleResponse, error) {
	vehicle, err := u.vehicleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("vehicle not found")
		}
		return nil, err
	}

	return vehicle.ToResponse(), nil
}

func (u *vehicleUsecase) GetAllVehicles(limit, offset int, search string) ([]entities.VehicleResponse, int64, error) {
	vehicles, total, err := u.vehicleRepo.GetAllVehicles(limit, offset, search)
	if err != nil {
		return nil, 0, errors.New("failed to fetch vehicles")
	}

	responses := make([]entities.VehicleResponse, len(vehicles))
	for i, vehicle := range vehicles {
		responses[i] = *vehicle.ToResponse()
	}

	return responses, total, nil
}

func (u *vehicleUsecase) UpdateVehicle(id uuid.UUID, req UpdateVehicleRequest) (*entities.VehicleResponse, error) {
	vehicle, err := u.vehicleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("vehicle not found")
		}
		return nil, err
	}

	if req.LicensePlate != "" {
		vehicle.LicensePlate = req.LicensePlate
	}
	if req.Model != "" {
		vehicle.Model = req.Model
	}

	if err := u.vehicleRepo.Update(vehicle); err != nil {
		return nil, errors.New("failed to update vehicle")
	}

	return vehicle.ToResponse(), nil
}

func (u *vehicleUsecase) DeleteVehicle(id uuid.UUID) error {
	_, err := u.vehicleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("vehicle not found")
		}
		return err
	}

	if err := u.vehicleRepo.Delete(id); err != nil {
		return errors.New("failed to delete vehicle")
	}

	return nil
}
