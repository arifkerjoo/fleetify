package main

import (
	"backend/internal/config"
	"backend/internal/infrastructure/database"
	"backend/internal/infrastructure/repositories"
	"backend/internal/interfaces/handlers"
	"backend/internal/interfaces/routes"
	"backend/internal/usecase"
	"backend/utils"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMysqlDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		panic(fmt.Sprintf("Migration failed: %v", err))
	}

	// Seeder
	if len(os.Args) > 1 {
		switch os.Args[1] {

		case "seed":
			if err := database.RunSeeders(db); err != nil {
				log.Fatalf("seeder failed: %v", err)
			}

			log.Println("Seeder success")
			return
		}
	}

	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret)
	userRepo := repositories.NewUserRepository(db)

	// Auth
	authUsecase := usecase.NewAuthUserUsecase(userRepo, jwtUtil)
	authHandler := handlers.NewAuthHandler(authUsecase, jwtUtil)

	// User
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUsecase, jwtUtil)

	//Vehicle
	vehicleRepo := repositories.NewVehicleRepository(db)
	vehicleUsecase := usecase.NewVehicleUsecase(vehicleRepo)
	vehicleHandler := handlers.NewVehicleHandler(vehicleUsecase, jwtUtil)

	//MasterItem
	masterItemRepo := repositories.NewMasterItemRepository(db)
	masterItemUsecase := usecase.NewMasterItemUsecase(masterItemRepo)
	masterItemHandler := handlers.NewMasterItemHandler(masterItemUsecase, jwtUtil)

	// MaintenanceReport
	maintenanceRepo := repositories.NewMaintenanceRepository(db)
	maintenanceUsecase := usecase.NewMaintenanceReportUsecase(maintenanceRepo, masterItemRepo, "")
	maintenanceHandler := handlers.NewMaintenanceReportHandler(maintenanceUsecase, jwtUtil)

	// Bundle handlers
	h := &routes.Handlers{
		Auth:              authHandler,
		User:              userHandler,
		Vehicle:           vehicleHandler,
		MasterItem:        masterItemHandler,
		MaintenanceReport: maintenanceHandler,
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   err.Error(),
				"code":    code,
				"message": "Something went wrong.",
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.AllowedOrigins, ","),
		AllowCredentials: true,
		AllowMethods: []string{
			"GET",
			"POST",
			"HEAD",
			"PUT",
			"DELETE",
			"PATCH",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
	}))

	app.Use("/uploads", static.New("./uploads", static.Config{
		Browse:     false,
		Compress:   true,
		ByteRange:  true,
		MaxAge:     3600,
		IndexNames: []string{},
		ModifyResponse: func(c fiber.Ctx) error {
			if c.Response().StatusCode() == 404 {
				return fiber.ErrNotFound
			}
			return nil
		},
	}))

	// Routes
	routes.SetupRoutes(app, h, jwtUtil)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
