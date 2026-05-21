package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportStatus string

const (
	StatusPendingApproval ReportStatus = "PENDING_APPROVAL"
	StatusApproved        ReportStatus = "APPROVED"
	StatusCompleted       ReportStatus = "COMPLETED"
	StatusRejected        ReportStatus = "REJECTED"
)

type MaintenanceReport struct {
	BaseModel
	VehicleID       uuid.UUID    `gorm:"not null;index"                                                                           json:"vehicle_id"`
	CreatedBy       uuid.UUID    `gorm:"not null;index"                                                                           json:"created_by"`
	AssignedTo      *uuid.UUID   `gorm:"index"                                                                                    json:"assigned_to"`
	Odometer        uint         `gorm:"type:int unsigned;not null"                                                               json:"odometer"`
	Complaint       string       `gorm:"type:text;not null"                                                                       json:"complaint"`
	Status          ReportStatus `gorm:"type:enum('PENDING_APPROVAL','APPROVED','COMPLETED','REJECTED');index;not null;default:'PENDING_APPROVAL'" json:"status"`
	InitialPhotoURL string       `gorm:"type:varchar(500)"                                                                        json:"initial_photo_url,omitempty"`
	ProofPhotoURL   string       `gorm:"type:varchar(500)"                                                                        json:"proof_photo_url,omitempty"`
	ApprovedAt      *time.Time   `gorm:"index"                                                                                    json:"approved_at,omitempty"`
	CompletedAt     *time.Time   `                                                                                                json:"completed_at,omitempty"`

	Vehicle     *Vehicle     `gorm:"foreignKey:VehicleID;references:ID"  json:"vehicle,omitempty"`
	Creator     *User        `gorm:"foreignKey:CreatedBy;references:ID"  json:"creator,omitempty"`
	ReportItems []ReportItem `gorm:"foreignKey:ReportID"                 json:"items,omitempty"`
}

func (r *MaintenanceReport) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (MaintenanceReport) TableName() string {
	return "maintenance_reports"
}

type ReportItem struct {
	BaseModel
	ReportID      uuid.UUID `gorm:"not null;index"                        json:"report_id"`
	ItemID        uuid.UUID `gorm:"not null;index"                        json:"item_id"`
	Quantity      int       `gorm:"type:int unsigned;not null;default:1"  json:"quantity"`
	PriceSnapshot float64   `gorm:"type:decimal(12,2);not null"           json:"price_snapshot"`
	Subtotal      float64   `gorm:"type:decimal(12,2);not null"           json:"subtotal"`

	Item *MasterItem `gorm:"foreignKey:ItemID;references:ID" json:"item,omitempty"`
}

func (r *ReportItem) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	// Hitung estimasi harga part/jasa
	r.Subtotal = r.PriceSnapshot * float64(r.Quantity)
	return nil
}

func (ReportItem) TableName() string {
	return "report_items"
}

type ReportItemResponse struct {
	ID            uuid.UUID           `json:"id"`
	ItemID        uuid.UUID           `json:"item_id"`
	Item          *MasterItemResponse `json:"item,omitempty"`
	Quantity      int                 `json:"quantity"`
	PriceSnapshot float64             `json:"price_snapshot"`
	Subtotal      float64             `json:"subtotal"`
}

type MaintenanceReportResponse struct {
	ID              uuid.UUID            `json:"id"`
	VehicleID       uuid.UUID            `json:"vehicle_id"`
	CreatedBy       uuid.UUID            `json:"created_by"`
	AssignedTo      *uuid.UUID           `json:"assigned_to,omitempty"`
	Odometer        uint                 `json:"odometer"`
	Complaint       string               `json:"complaint"`
	Status          ReportStatus         `json:"status"`
	InitialPhotoURL string               `json:"initial_photo_url,omitempty"`
	ProofPhotoURL   string               `json:"proof_photo_url,omitempty"`
	ApprovedAt      *time.Time           `json:"approved_at,omitempty"`
	CompletedAt     *time.Time           `json:"completed_at,omitempty"`
	TotalAmount     float64              `json:"total_amount"`
	Vehicle         *VehicleResponse     `json:"vehicle,omitempty"`
	Creator         *UserResponse        `json:"creator,omitempty"`
	Items           []ReportItemResponse `json:"items,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}
