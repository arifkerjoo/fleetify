package routes

import (
	"backend/internal/domain/entities"
	"backend/internal/interfaces/handlers"
	"backend/internal/interfaces/middleware"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

func MaintenanceReportRoutes(r fiber.Router, h *handlers.MaintenanceReportHandler, jwtUtil *utils.JWTUtil) {
	reports := r.Group("/reports")
	reports.Use(middleware.JWTMiddleware(jwtUtil))

	setupReportReadRoutes(reports, h)
	setupReportSARoutes(reports, h)
	setupReportApprovalRoutes(reports, h)
}

func setupReportReadRoutes(r fiber.Router, h *handlers.MaintenanceReportHandler) {
	r.Get("", h.GetAllReports)
	r.Get("/:id", h.GetReportByID)
}

// hanya role SA
func setupReportSARoutes(r fiber.Router, h *handlers.MaintenanceReportHandler) {
	sa := r.Group("")
	sa.Use(middleware.RequireRoleMiddleware(string(entities.RoleSA)))

	// F-01: buat laporan baru
	sa.Post("", h.CreateReport)

	// F-01 supplemental: upload foto awal
	sa.Post("/:id/initial-photo", h.UploadInitialPhoto)

	// F-03: selesaikan pengerjaan dengan foto bukti
	sa.Patch("/:id/complete", h.CompleteReport)
}

// hanya role APPROVAL
func setupReportApprovalRoutes(r fiber.Router, h *handlers.MaintenanceReportHandler) {
	approval := r.Group("")
	approval.Use(middleware.RequireRoleMiddleware(string(entities.RoleApproval)))

	// F-02: setujui laporan
	approval.Patch("/:id/approve", h.ApproveReport)
}
