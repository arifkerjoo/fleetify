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
	Odometer        uuid.UUID    `gorm:"type:int unsigned;not null"                                                               json:"odometer"`
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
	ID            uuid.UUID `gorm:"primaryKey;autoIncrement"              json:"-"`
	UUID          string    `gorm:"type:varchar(36);uniqueIndex;not null" json:"id"`
	ReportID      uuid.UUID `gorm:"not null;index"                        json:"report_id"`
	ItemID        uuid.UUID `gorm:"not null;index"                        json:"item_id"`
	Quantity      int       `gorm:"type:int unsigned;not null;default:1"  json:"quantity"`
	PriceSnapshot float64   `gorm:"type:decimal(12,2);not null"           json:"price_snapshot"`
	Subtotal      float64   `gorm:"type:decimal(12,2);not null"           json:"subtotal"`
	CreatedAt     time.Time `gorm:"autoCreateTime"                        json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"                        json:"updated_at"`

	Item *MasterItem `gorm:"foreignKey:ItemID;references:ID" json:"item,omitempty"`
}

func (r *ReportItem) BeforeCreate(tx *gorm.DB) error {
	if r.UUID == "" {
		r.UUID = uuid.New().String()
	}
	// Hitung estimasi harga part/jasa
	r.Subtotal = r.PriceSnapshot * float64(r.Quantity)
	return nil
}

func (ReportItem) TableName() string {
	return "report_items"
}
