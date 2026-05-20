package database

import (
	"backend/internal/seeder"

	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) error {
	if err := seeder.SeedUsers(db); err != nil {
		return err
	}

	if err := seeder.SeedVehicles(db); err != nil {
		return err
	}

	if err := seeder.SeedMasterItems(db); err != nil {
		return err
	}

	return nil
}
