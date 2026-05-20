package seeder

import (
	"backend/internal/domain/entities"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) error {
	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedVehicles(db); err != nil {
		return err
	}

	if err := SeedMasterItems(db); err != nil {
		return err
	}

	return nil
}

func SeedUsers(db *gorm.DB) error {
	var count int64

	db.Model(&entities.User{}).Count(&count)

	if count > 0 {
		log.Println("users already seeded")
		return nil
	}

	now := time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	users := []entities.User{
		{
			Username:   "superadmin",
			Email:      "superadmin@fleetify.com",
			Password:   string(hashedPassword),
			FullName:   "Super Admin",
			Phone:      "081111111111",
			Role:       entities.UserRole("SA"),
			IsActive:   true,
			VerifiedAt: &now,
		},
		{
			Username:   "approval1",
			Email:      "approval1@fleetify.com",
			Password:   string(hashedPassword),
			FullName:   "Approval User 1",
			Phone:      "082222222222",
			Role:       entities.UserRole("APPROVAL"),
			IsActive:   true,
			VerifiedAt: &now,
		},
		{
			Username:   "approval2",
			Email:      "approval2@fleetify.com",
			Password:   string(hashedPassword),
			FullName:   "Approval User 2",
			Phone:      "083333333333",
			Role:       entities.UserRole("APPROVAL"),
			IsActive:   true,
			VerifiedAt: &now,
		},
	}

	for _, user := range users {
		u := user

		if err := db.
			Where("email = ?", u.Email).
			FirstOrCreate(&u).Error; err != nil {
			return err
		}

		log.Printf("seeded user: %s", u.Email)
	}

	return nil
}

func SeedVehicles(db *gorm.DB) error {
	var count int64

	db.Model(&entities.Vehicle{}).Count(&count)

	if count > 0 {
		log.Println("vehicles already seeded")
		return nil
	}

	vehicles := []entities.Vehicle{
		{
			LicensePlate: "B1234AAA",
			Brand:        "Toyota",
			Model:        "Avanza",
			Year:         2021,
			Status:       entities.VehicleStatusActive,
		},
		{
			LicensePlate: "B2345BBB",
			Brand:        "Honda",
			Model:        "Brio",
			Year:         2022,
			Status:       entities.VehicleStatusActive,
		},
		{
			LicensePlate: "B3456CCC",
			Brand:        "Suzuki",
			Model:        "Ertiga",
			Year:         2020,
			Status:       entities.VehicleStatusActive,
		},
		{
			LicensePlate: "B4567DDD",
			Brand:        "Mitsubishi",
			Model:        "L300",
			Year:         2019,
			Status:       entities.VehicleStatusInactive,
		},
		{
			LicensePlate: "B5678EEE",
			Brand:        "Daihatsu",
			Model:        "Gran Max",
			Year:         2023,
			Status:       entities.VehicleStatusActive,
		},
	}

	for _, vehicle := range vehicles {
		v := vehicle

		if err := db.
			Where("license_plate = ?", v.LicensePlate).
			FirstOrCreate(&v).Error; err != nil {
			return err
		}

		log.Printf("seeded vehicle: %s", v.LicensePlate)
	}

	return nil
}

func SeedMasterItems(db *gorm.DB) error {
	var count int64

	db.Model(&entities.MasterItem{}).Count(&count)

	if count > 0 {
		log.Println("master items already seeded")
		return nil
	}

	items := []entities.MasterItem{
		{
			ItemName: "Engine Oil",
			ItemType: entities.ItemTypePart,
			Price:    350000,
			IsActive: true,
		},
		{
			ItemName: "Brake Pad",
			ItemType: entities.ItemTypePart,
			Price:    450000,
			IsActive: true,
		},
		{
			ItemName: "Oil Filter",
			ItemType: entities.ItemTypePart,
			Price:    120000,
			IsActive: true,
		},
		{
			ItemName: "Air Filter",
			ItemType: entities.ItemTypePart,
			Price:    150000,
			IsActive: true,
		},
		{
			ItemName: "Battery",
			ItemType: entities.ItemTypePart,
			Price:    950000,
			IsActive: true,
		},
		{
			ItemName: "General Service",
			ItemType: entities.ItemTypeService,
			Price:    250000,
			IsActive: true,
		},
		{
			ItemName: "Engine Tune Up",
			ItemType: entities.ItemTypeService,
			Price:    500000,
			IsActive: true,
		},
	}

	for _, item := range items {
		i := item

		if err := db.
			Where("item_name = ?", i.ItemName).
			FirstOrCreate(&i).Error; err != nil {
			return err
		}

		log.Printf("seeded master item: %s", i.ItemName)
	}

	return nil
}
