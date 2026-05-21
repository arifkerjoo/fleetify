package routes

import (
	"backend/internal/interfaces/handlers"
	"backend/internal/interfaces/middleware"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

func VehicleRoutes(r fiber.Router, vehicleHandler *handlers.VehicleHandler, jwtUtil *utils.JWTUtil) {
	vehicles := r.Group("/vehicles")

	vehicles.Use(middleware.JWTMiddleware(jwtUtil))

	setupVehicleReadRoutes(vehicles, vehicleHandler)

	setupVehicleWriteRoutes(vehicles, vehicleHandler)
}

func setupVehicleReadRoutes(r fiber.Router, h *handlers.VehicleHandler) {
	r.Get("", h.GetAllVehicles)
	r.Get("/:id", h.GetVehicleByID)
}

func setupVehicleWriteRoutes(r fiber.Router, h *handlers.VehicleHandler) {
	approval := r.Group("")
	approval.Use(middleware.RequireRoleMiddleware("APPROVAL"))

	approval.Post("", h.CreateVehicle)
	approval.Put("/:id", h.UpdateVehicle)
	approval.Delete("/:id", h.DeleteVehicle)
}
