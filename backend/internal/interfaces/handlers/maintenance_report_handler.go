package handlers

import (
	"backend/internal/domain/entities"
	"backend/internal/interfaces/middleware"
	"backend/internal/usecase"
	"backend/utils"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MaintenanceReportHandler struct {
	reportUsecase usecase.MaintenanceReportUsecase
	jwtUtil       *utils.JWTUtil
}

func NewMaintenanceReportHandler(
	reportUsecase usecase.MaintenanceReportUsecase,
	jwtUtil *utils.JWTUtil,
) *MaintenanceReportHandler {
	return &MaintenanceReportHandler{
		reportUsecase: reportUsecase,
		jwtUtil:       jwtUtil,
	}
}

func (h *MaintenanceReportHandler) RegisterRoutes(api fiber.Router) {
	reports := api.Group("/reports")
	reports.Use(middleware.JWTMiddleware(h.jwtUtil))

	reports.Get("", h.GetAllReports)
	reports.Get("/:id", h.GetReportByID)

	// F-01: SA
	reports.Post("",
		middleware.RequireRoleMiddleware(string(entities.RoleSA)),
		h.CreateReport,
	)

	// F-01 upload photo
	reports.Post("/:id/initial-photo",
		middleware.RequireRoleMiddleware(string(entities.RoleSA)),
		h.UploadInitialPhoto,
	)

	// F-02: Approval
	reports.Patch("/:id/approve",
		middleware.RequireRoleMiddleware(string(entities.RoleApproval)),
		h.ApproveReport,
	)

	// F-03: SA complete
	reports.Patch("/:id/complete",
		middleware.RequireRoleMiddleware(string(entities.RoleSA)),
		h.CompleteReport,
	)
}

func (h *MaintenanceReportHandler) GetAllReports(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	search := c.Query("search", "")
	status := c.Query("status", "")

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	if status != "" {
		validStatuses := map[string]bool{
			string(entities.StatusPendingApproval): true,
			string(entities.StatusApproved):        true,
			string(entities.StatusCompleted):       true,
			string(entities.StatusRejected):        true,
		}
		if !validStatuses[status] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Status Invalid. Yang diperbolehkan: PENDING_APPROVAL, APPROVED, COMPLETED, REJECTED",
			})
		}
	}

	offset := (page - 1) * limit

	reports, total, err := h.reportUsecase.GetAllReports(limit, offset, search, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    reports,
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

func (h *MaintenanceReportHandler) GetReportByID(c fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid report ID",
		})
	}

	report, err := h.reportUsecase.GetReportByID(reportID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    report,
	})
}

func (h *MaintenanceReportHandler) CreateReport(c fiber.Ctx) error {
	createdBy, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	var req usecase.CreateReportRequest
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

	report, err := h.reportUsecase.CreateReport(createdBy, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Maintenance report created successfully",
		"data":    report,
	})
}

func (h *MaintenanceReportHandler) UploadInitialPhoto(c fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid report ID",
		})
	}

	requesterID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	existing, err := h.reportUsecase.GetReportByID(reportID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if existing.CreatedBy != requesterID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "You can only upload photos to your own reports",
		})
	}

	if existing.Status != entities.StatusPendingApproval {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Initial photo can only be uploaded for PENDING_APPROVAL reports",
		})
	}

	photoURL, err := h.savePhoto(c, "initial_photo", "reports/initial", reportID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	updateReq := usecase.CreateReportRequest{
		VehicleID:       existing.VehicleID,
		Odometer:        existing.Odometer,
		Complaint:       existing.Complaint,
		InitialPhotoURL: photoURL,
		Items:           nil,
	}
	_ = updateReq

	patchReq := usecase.PatchInitialPhotoRequest{PhotoURL: photoURL}
	report, err := h.reportUsecase.PatchInitialPhoto(reportID, patchReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Foto berhasil diupload",
		"data":    report,
	})
}

func (h *MaintenanceReportHandler) ApproveReport(c fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid report ID format",
		})
	}

	approverID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	var req usecase.ApproveReportRequest
	_ = c.Bind().Body(&req)

	report, err := h.reportUsecase.ApproveReport(reportID, approverID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Report berhasil di approve",
		"data":    report,
	})
}

func (h *MaintenanceReportHandler) CompleteReport(c fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid report ID",
		})
	}

	saID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	photoURL, err := h.savePhoto(c, "proof_photo", "reports/proof", reportID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	req := usecase.CompleteReportRequest{ProofPhotoURL: photoURL}

	report, err := h.reportUsecase.CompleteReport(reportID, saID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Status berhasil diubah",
		"data":    report,
	})
}

func (h *MaintenanceReportHandler) savePhoto(
	c fiber.Ctx,
	formKey, subDir, prefix string,
) (string, error) {

	file, err := c.FormFile(formKey)
	if err != nil {
		return "", fmt.Errorf("photo field '%s' is required", formKey)
	}

	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("ukuran file tidak boleh lebih dari 5MB")
	}

	allowedExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExt[ext] {
		return "", fmt.Errorf("hanya boleh format jpg, jpeg, png, webp")
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", prefix, timestamp, ext)
	uploadPath := filepath.Join("uploads", subDir, filename)

	if err := c.SaveFile(file, uploadPath); err != nil {
		return "", fmt.Errorf("Gagal menyimpan foto")
	}

	return fmt.Sprintf("/uploads/%s/%s", subDir, filename), nil
}
