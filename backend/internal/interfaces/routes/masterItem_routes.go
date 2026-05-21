package routes

import (
	"backend/internal/interfaces/handlers"
	"backend/internal/interfaces/middleware"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

func MasterItemRoutes(r fiber.Router, masterItemHandler *handlers.MasterItemHandler, jwtUtil *utils.JWTUtil) {
	items := r.Group("/master-items")

	items.Use(middleware.JWTMiddleware(jwtUtil))

	setupMasterItemReadRoutes(items, masterItemHandler)

	setupMasterItemWriteRoutes(items, masterItemHandler)
}

func setupMasterItemReadRoutes(r fiber.Router, h *handlers.MasterItemHandler) {
	r.Get("", h.GetAllMasterItems)
	r.Get("/:id", h.GetMasterItemByID)
}

func setupMasterItemWriteRoutes(r fiber.Router, h *handlers.MasterItemHandler) {
	approval := r.Group("")
	approval.Use(middleware.RequireRoleMiddleware("approval"))

	approval.Post("", h.CreateMasterItem)
	approval.Put("/:id", h.UpdateMasterItem)
	approval.Delete("/:id", h.DeleteMasterItem)
}
