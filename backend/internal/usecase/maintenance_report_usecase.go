package usecase

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repositories"
	"backend/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportItemRequest struct {
	ItemID   uuid.UUID `json:"item_id"   validate:"required"`
	Quantity int       `json:"quantity"  validate:"required,min=1"`
}

type CreateReportRequest struct {
	VehicleID       uuid.UUID           `json:"vehicle_id"        validate:"required"`
	Odometer        uint                `json:"odometer"          validate:"required,min=1"`
	Complaint       string              `json:"complaint"         validate:"required"`
	InitialPhotoURL string              `json:"initial_photo_url"`
	Items           []ReportItemRequest `json:"items"             validate:"required,min=1,dive"`
}

// F-02 (APPROVAL)
type ApproveReportRequest struct {
	AssignedTo *uuid.UUID `json:"assigned_to"`
}

// F-03 (SA)
type CompleteReportRequest struct {
	ProofPhotoURL string `json:"proof_photo_url" validate:"required"`
}

type MaintenanceReportUsecase interface {
	CreateReport(createdBy uuid.UUID, req CreateReportRequest) (*entities.MaintenanceReportResponse, error)                        // F-01: SA create report [PENDING_APPROVAL]
	ApproveReport(reportID uuid.UUID, approverID uuid.UUID, req ApproveReportRequest) (*entities.MaintenanceReportResponse, error) // F-02: Approval [APPROVEED]
	CompleteReport(reportID uuid.UUID, saID uuid.UUID, req CompleteReportRequest) (*entities.MaintenanceReportResponse, error)     // F-03: SA upload photo
	GetAllReports(limit, offset int, search, status string) ([]entities.MaintenanceReportResponse, int64, error)                   // F-04: List all reports
	GetReportByID(id uuid.UUID) (*entities.MaintenanceReportResponse, error)

	PatchInitialPhoto(reportID uuid.UUID, req PatchInitialPhotoRequest) (*entities.MaintenanceReportResponse, error)
}

type maintenanceReportUsecase struct {
	reportRepo     repositories.MaintenanceRepository
	masterItemRepo repositories.MasterItemRepository
	webhookURL     string // add B-02
}

func NewMaintenanceReportUsecase(
	reportRepo repositories.MaintenanceRepository,
	masterItemRepo repositories.MasterItemRepository,
	webhookURL string,
) MaintenanceReportUsecase {
	return &maintenanceReportUsecase{
		reportRepo:     reportRepo,
		masterItemRepo: masterItemRepo,
		webhookURL:     webhookURL,
	}
}

func (u *maintenanceReportUsecase) CreateReport(
	createdBy uuid.UUID,
	req CreateReportRequest,
) (*entities.MaintenanceReportResponse, error) {

	reportItems, err := u.buildReportItems(uuid.Nil, req.Items)
	if err != nil {
		return nil, err
	}

	report := &entities.MaintenanceReport{
		VehicleID:       req.VehicleID,
		CreatedBy:       createdBy,
		Complaint:       req.Complaint,
		InitialPhotoURL: req.InitialPhotoURL,
		Status:          entities.StatusPendingApproval,
		Odometer:        req.Odometer,
	}

	if err := u.reportRepo.CreateWithItems(report, reportItems); err != nil {
		return nil, errors.New("Gagal create maintenance report")
	}

	return u.GetReportByID(report.ID)
}

// ── F-02 ─────────────────────────────────────

func (u *maintenanceReportUsecase) ApproveReport(
	reportID uuid.UUID,
	approverID uuid.UUID,
	req ApproveReportRequest,
) (*entities.MaintenanceReportResponse, error) {

	report, err := u.reportRepo.GetByID(reportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report not found")
		}
		return nil, err
	}

	if report.Status != entities.StatusPendingApproval {
		return nil, errors.New("only PENDING_APPROVAL reports can be approved")
	}

	now := time.Now()
	report.Status = entities.StatusApproved
	report.ApprovedAt = &now

	if req.AssignedTo != nil {
		report.AssignedTo = req.AssignedTo
	}

	if err := u.reportRepo.Update(report); err != nil {
		return nil, errors.New("Gagal approve report")
	}

	// B-02: webhook
	utils.SendWebhook(u.webhookURL, utils.Payload{
		Event:     "report.approved",
		ReportID:  report.ID.String(),
		Status:    string(entities.StatusApproved),
		Timestamp: now,
		Data: map[string]interface{}{
			"approved_by": approverID.String(),
			"assigned_to": req.AssignedTo,
		},
	})

	return u.GetReportByID(report.ID)
}

// ── F-03 ─────────────────────────────────────

func (u *maintenanceReportUsecase) CompleteReport(
	reportID uuid.UUID,
	saID uuid.UUID,
	req CompleteReportRequest,
) (*entities.MaintenanceReportResponse, error) {

	report, err := u.reportRepo.GetByID(reportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report not found")
		}
		return nil, err
	}

	if report.Status != entities.StatusApproved {
		return nil, errors.New("hanya report dengan status APPROVED yg dapat di update")
	}

	if report.AssignedTo != nil && *report.AssignedTo != saID {
		return nil, errors.New("Gagal mengubah status")
	}

	now := time.Now()
	report.Status = entities.StatusCompleted
	report.CompletedAt = &now
	report.ProofPhotoURL = req.ProofPhotoURL

	if err := u.reportRepo.Update(report); err != nil {
		return nil, errors.New("gagal complete report")
	}

	// B-02: webhook
	utils.SendWebhook(u.webhookURL, utils.Payload{
		Event:     "report.completed",
		ReportID:  report.ID.String(),
		Status:    string(entities.StatusCompleted),
		Timestamp: now,
		Data: map[string]interface{}{
			"completed_by": saID.String(),
		},
	})

	return u.GetReportByID(report.ID)
}

// ── F-04 ─────────────────────────────────────

func (u *maintenanceReportUsecase) GetAllReports(
	limit, offset int,
	search, status string,
) ([]entities.MaintenanceReportResponse, int64, error) {

	reports, total, err := u.reportRepo.GetAll(limit, offset, search, status)
	if err != nil {
		return nil, 0, errors.New("gagal fetch reports")
	}

	responses := make([]entities.MaintenanceReportResponse, len(reports))
	for i, r := range reports {
		responses[i] = toReportResponse(&r)
	}

	return responses, total, nil
}

func (u *maintenanceReportUsecase) GetReportByID(id uuid.UUID) (*entities.MaintenanceReportResponse, error) {
	report, err := u.reportRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report not found")
		}
		return nil, err
	}

	res := toReportResponse(report)
	return &res, nil
}

func (u *maintenanceReportUsecase) buildReportItems(
	reportID uuid.UUID,
	reqs []ReportItemRequest,
) ([]entities.ReportItem, error) {

	items := make([]entities.ReportItem, 0, len(reqs))

	for _, r := range reqs {
		masterItem, err := u.masterItemRepo.GetByID(r.ItemID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("item not found: " + r.ItemID.String())
			}
			return nil, err
		}

		if !masterItem.IsActive {
			return nil, errors.New("item is no longer active: " + masterItem.ItemName)
		}

		items = append(items, entities.ReportItem{
			ReportID:      reportID,
			ItemID:        r.ItemID,
			Quantity:      r.Quantity,
			PriceSnapshot: masterItem.Price,
			Subtotal:      masterItem.Price * float64(r.Quantity),
		})
	}

	return items, nil
}

func toReportResponse(r *entities.MaintenanceReport) entities.MaintenanceReportResponse {
	var totalAmount float64
	items := make([]entities.ReportItemResponse, 0, len(r.ReportItems))

	for _, ri := range r.ReportItems {
		totalAmount += ri.Subtotal

		var itemResp *entities.MasterItemResponse
		if ri.Item != nil {
			itemResp = ri.Item.ToResponse()
		}

		items = append(items, entities.ReportItemResponse{
			ID:            ri.ID,
			ItemID:        ri.ItemID,
			Item:          itemResp,
			Quantity:      ri.Quantity,
			PriceSnapshot: ri.PriceSnapshot,
			Subtotal:      ri.Subtotal,
		})
	}

	var vehicleResp *entities.VehicleResponse
	if r.Vehicle != nil {
		vehicleResp = r.Vehicle.ToResponse()
	}

	var creatorResp *entities.UserResponse
	if r.Creator != nil {
		creatorResp = r.Creator.ToResponse()
	}

	return entities.MaintenanceReportResponse{
		ID:              r.ID,
		VehicleID:       r.VehicleID,
		CreatedBy:       r.CreatedBy,
		AssignedTo:      r.AssignedTo,
		Complaint:       r.Complaint,
		Status:          r.Status,
		InitialPhotoURL: r.InitialPhotoURL,
		ProofPhotoURL:   r.ProofPhotoURL,
		ApprovedAt:      r.ApprovedAt,
		CompletedAt:     r.CompletedAt,
		TotalAmount:     totalAmount,
		Vehicle:         vehicleResp,
		Creator:         creatorResp,
		Items:           items,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

type PatchInitialPhotoRequest struct {
	PhotoURL string `json:"photo_url" validate:"required"`
}

func (u *maintenanceReportUsecase) PatchInitialPhoto(
	reportID uuid.UUID,
	req PatchInitialPhotoRequest,
) (*entities.MaintenanceReportResponse, error) {

	report, err := u.reportRepo.GetByID(reportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report tidak ditemukan")
		}
		return nil, err
	}

	if report.Status != entities.StatusPendingApproval {
		return nil, errors.New("Foto hanya dapat di update ketika statsus PENDING_APPROVAL")
	}

	report.InitialPhotoURL = req.PhotoURL

	if err := u.reportRepo.Update(report); err != nil {
		return nil, errors.New("gagal update photo")
	}

	return u.GetReportByID(report.ID)
}
