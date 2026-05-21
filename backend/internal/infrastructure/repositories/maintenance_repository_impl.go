package repositories

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type maintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) repositories.MaintenanceRepository {
	return &maintenanceRepository{db: db}
}

func (r *maintenanceRepository) CreateWithItems(
	report *entities.MaintenanceReport,
	items []entities.ReportItem,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(report).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].ReportID = report.ID
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *maintenanceRepository) GetByID(id uuid.UUID) (*entities.MaintenanceReport, error) {
	var report entities.MaintenanceReport
	err := r.db.
		Preload("Vehicle").
		Preload("Creator").
		Preload("ReportItems").
		Preload("ReportItems.Item").
		Where("id = ?", id).
		First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *maintenanceRepository) Update(report *entities.MaintenanceReport) error {
	return r.db.Save(report).Error
}

func (r *maintenanceRepository) GetAll(
	limit, offset int,
	search, status string,
) ([]entities.MaintenanceReport, int64, error) {

	var reports []entities.MaintenanceReport
	var total int64

	query := r.db.Model(&entities.MaintenanceReport{}).
		Joins("LEFT JOIN vehicles ON vehicles.id = maintenance_reports.vehicle_id").
		Joins("LEFT JOIN users   ON users.id   = maintenance_reports.created_by").
		Preload("Vehicle").
		Preload("Creator").
		Preload("ReportItems").
		Preload("ReportItems.Item")

	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"vehicles.license_plate LIKE ? OR users.full_name LIKE ?",
			like, like,
		)
	}

	if status != "" {
		query = query.Where("maintenance_reports.status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("maintenance_reports.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *maintenanceRepository) DeleteReportItemsByReportID(reportID uuid.UUID) error {
	return r.db.
		Where("report_id = ?", reportID).
		Delete(&entities.ReportItem{}).Error
}
