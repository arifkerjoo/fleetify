package handlers

import (
	"backend/internal/interfaces/middleware"
	"backend/internal/usecase"
	"backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MasterItemHandler struct {
	masterItemUsecase usecase.MasterItemUsecase
	jwtUtil           *utils.JWTUtil
}

func NewMasterItemHandler(masterItemUsecase usecase.MasterItemUsecase, jwtUtil *utils.JWTUtil) *MasterItemHandler {
	return &MasterItemHandler{
		masterItemUsecase: masterItemUsecase,
		jwtUtil:           jwtUtil,
	}
}

func (h *MasterItemHandler) RegisterRoutes(api fiber.Router) {
	items := api.Group("/master-items")
	items.Use(middleware.JWTMiddleware(h.jwtUtil))

	items.Get("", h.GetAllMasterItems)
	items.Get("/:id", h.GetMasterItemByID)

	items.Post("", middleware.RequireRoleMiddleware("admin"), h.CreateMasterItem)
	items.Put("/:id", middleware.RequireRoleMiddleware("admin"), h.UpdateMasterItem)
	items.Delete("/:id", middleware.RequireRoleMiddleware("admin"), h.DeleteMasterItem)
}

func (h *MasterItemHandler) GetAllMasterItems(c fiber.Ctx) error {
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

	items, total, err := h.masterItemUsecase.GetAllMasterItems(limit, offset, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    items,
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

func (h *MasterItemHandler) GetMasterItemByID(c fiber.Ctx) error {
	itemID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid master item ID format",
		})
	}

	item, err := h.masterItemUsecase.GetMasterItemByID(itemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    item,
	})
}

func (h *MasterItemHandler) CreateMasterItem(c fiber.Ctx) error {
	var req usecase.CreateMasterItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	item, err := h.masterItemUsecase.CreateMasterItem(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Master item created successfully",
		"data":    item,
	})
}

func (h *MasterItemHandler) UpdateMasterItem(c fiber.Ctx) error {
	itemID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid master item ID format",
		})
	}

	var req usecase.UpdateMasterItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	item, err := h.masterItemUsecase.UpdateMasterItem(itemID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Master item updated successfully",
		"data":    item,
	})
}

func (h *MasterItemHandler) DeleteMasterItem(c fiber.Ctx) error {
	itemID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid master item ID format",
		})
	}

	if err := h.masterItemUsecase.DeleteMasterItem(itemID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Master item deleted successfully",
	})
}
