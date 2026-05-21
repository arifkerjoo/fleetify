package repositories

import (
	"backend/internal/domain/entities"

	"github.com/google/uuid"
)

type MaintenanceRepository interface {
	CreateWithItems(report *entities.MaintenanceReport, items []entities.ReportItem) error
	GetByID(id uuid.UUID) (*entities.MaintenanceReport, error)
	Update(report *entities.MaintenanceReport) error
	GetAll(limit, offset int, search, status string) ([]entities.MaintenanceReport, int64, error)
	DeleteReportItemsByReportID(reportID uuid.UUID) error
}
