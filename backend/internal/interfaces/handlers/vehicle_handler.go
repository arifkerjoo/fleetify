package handlers

import (
	"backend/internal/interfaces/middleware"
	"backend/internal/usecase"
	"backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type VehicleHandler struct {
	vehicleUsecase usecase.VehicleUsecase
	jwtUtil        *utils.JWTUtil
}

func NewVehicleHandler(vehicleUsecase usecase.VehicleUsecase, jwtUtil *utils.JWTUtil) *VehicleHandler {
	return &VehicleHandler{
		vehicleUsecase: vehicleUsecase,
		jwtUtil:        jwtUtil,
	}
}

func (h *VehicleHandler) RegisterRoutes(api fiber.Router) {
	vehicles := api.Group("/vehicles")
	vehicles.Use(middleware.JWTMiddleware(h.jwtUtil))

	vehicles.Get("", h.GetAllVehicles)
	vehicles.Get("/:id", h.GetVehicleByID)

	// approval role only
	vehicles.Post("", middleware.RequireRoleMiddleware("approval"), h.CreateVehicle)
	vehicles.Put("/:id", middleware.RequireRoleMiddleware("approval"), h.UpdateVehicle)
	vehicles.Delete("/:id", middleware.RequireRoleMiddleware("approval"), h.DeleteVehicle)
}

func (h *VehicleHandler) GetAllVehicles(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	search := c.Query("search", "")

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	vehicles, total, err := h.vehicleUsecase.GetAllVehicles(limit, offset, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    vehicles,
		"pagination": fiber.Map{
			"total":        total,
			"page":         page,
			"limit":        limit,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_previous": page > 1,
		},
	})
}

func (h *VehicleHandler) GetVehicleByID(c fiber.Ctx) error {
	vehicleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vehicle ID format",
		})
	}

	vehicle, err := h.vehicleUsecase.GetVehicleByID(vehicleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    vehicle,
	})
}

func (h *VehicleHandler) CreateVehicle(c fiber.Ctx) error {
	var req usecase.CreateVehicleRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	vehicle, err := h.vehicleUsecase.CreateVehicle(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Vehicle created successfully",
		"data":    vehicle,
	})
}

func (h *VehicleHandler) UpdateVehicle(c fiber.Ctx) error {
	vehicleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vehicle ID format",
		})
	}

	var req usecase.UpdateVehicleRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	vehicle, err := h.vehicleUsecase.UpdateVehicle(vehicleID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vehicle updated successfully",
		"data":    vehicle,
	})
}

func (h *VehicleHandler) DeleteVehicle(c fiber.Ctx) error {
	vehicleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vehicle ID format",
		})
	}

	if err := h.vehicleUsecase.DeleteVehicle(vehicleID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vehicle deleted successfully",
	})
}
