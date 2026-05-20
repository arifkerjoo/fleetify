package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleStatus string

const (
	VehicleStatusActive   VehicleStatus = "ACTIVE"
	VehicleStatusInactive VehicleStatus = "INACTIVE"
)

type Vehicle struct {
	BaseModel

	LicensePlate string        `gorm:"type:varchar(20);uniqueIndex;not null"             json:"license_plate"`
	Brand        string        `gorm:"type:varchar(100);not null"                        json:"brand"`
	Model        string        `gorm:"type:varchar(100);not null"                        json:"model"`
	Year         int           `gorm:"type:smallint unsigned;not null"                   json:"year"`
	Status       VehicleStatus `gorm:"type:enum('ACTIVE','INACTIVE');default:'ACTIVE';not null" json:"status"`

	MaintenanceReports []MaintenanceReport `gorm:"foreignKey:VehicleID" json:"-"`
}

func (v *Vehicle) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

func (Vehicle) TableName() string {
	return "vehicles"
}
