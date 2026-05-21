package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemType string

const (
	ItemTypePart    ItemType = "PART"
	ItemTypeService ItemType = "SERVICE"
)

type MasterItem struct {
	BaseModel

	ItemName string   `gorm:"type:varchar(255);not null"                 json:"item_name"`
	ItemType ItemType `gorm:"type:enum('PART','SERVICE');index;not null" json:"item_type"`
	Price    float64  `gorm:"type:decimal(12,2);not null"                json:"price"`
	IsActive bool     `gorm:"default:true;index;not null"                json:"is_active"`
}

type MasterItemResponse struct {
	ID       uuid.UUID `json:"id"`
	ItemName string    `json:"item_name"`
	ItemType ItemType  `json:"item_type"`
	Price    float64   `json:"price"`
	IsActive bool      `json:"is_active"`
}

func (m *MasterItem) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	return nil
}

func (MasterItem) TableName() string {
	return "master_items"
}

func (m *MasterItem) ToResponse() *MasterItemResponse {
	return &MasterItemResponse{
		ItemName: m.ItemName,
		ItemType: m.ItemType,
		Price:    m.Price,
		IsActive: m.IsActive,
	}
}
