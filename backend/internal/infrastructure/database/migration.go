package database

import (
	"backend/internal/domain/entities"

	"gorm.io/gorm"
)

func GetAllModels() []interface{} {
	return []interface{}{
		&entities.User{},
		&entities.Vehicle{},
		&entities.MasterItem{},
		&entities.MaintenanceReport{},
		&entities.ReportItem{},
	}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(GetAllModels()...)
}
