package repositories

import (
	"backend/internal/shared"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// Create
func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// Update
func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete
func (r *BaseRepository[T]) Delete(id string) error {
	var model T
	return r.db.Delete(&model, "id = ?", id).Error
}

// FindOne — ambil satu record berdasarkan conditions
func (r *BaseRepository[T]) FindOne(
	conditions map[string]interface{},
	options ...shared.FindOption,
) (*T, error) {
	var entity T
	query := r.db.Model(&entity)

	for _, option := range options {
		query = option.Apply(query)
	}

	for field, value := range conditions {
		if field == "" {
			continue
		}
		if strings.Contains(field, "?") {
			query = query.Where(field, value)
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	if err := query.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll — ambil banyak record berdasarkan kondisi
func (r *BaseRepository[T]) FindAll(
	conditions map[string]interface{},
	options ...shared.FindOption,
) ([]T, error) {
	var entities []T
	query := r.db.Model(new(T))

	for _, option := range options {
		query = option.Apply(query)
	}

	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// FindAllPaginated — dengan pagination dan search (opsional)
func (r *BaseRepository[T]) FindAllPaginated(
	limit, offset int,
	search string,
	searchFields []string,
) ([]T, int64, error) {
	var entities []T
	var total int64

	model := new(T)
	countQuery := r.db.Model(model)

	if search != "" && len(searchFields) > 0 {
		searchPattern := "%" + search + "%"
		orClauses := make([]string, 0, len(searchFields)) // Bikin slice kosong sebanyak search field nya untuk nama" field nya
		args := make([]interface{}, 0, len(searchFields)) // Bikin slice kosong sebanyak search field nya untuk value searchnya

		for _, field := range searchFields {
			orClauses = append(orClauses, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, searchPattern)
		}
		countQuery = countQuery.Where(strings.Join(orClauses, " OR "), args...)
	}

	countQuery.Count(&total)

	query := r.db.Model(model).Limit(limit).Offset(offset)
	if search != "" && len(searchFields) > 0 {
		searchPattern := "%" + search + "%"
		orClauses := make([]string, 0, len(searchFields)) // Bikin slice kosong sebanyak search field nya untuk nama" field nya
		args := make([]interface{}, 0, len(searchFields)) // Bikin slice kosong sebanyak search field nya untuk value searchnya

		for _, field := range searchFields {
			orClauses = append(orClauses, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, searchPattern)
		}
		query = query.Where(strings.Join(orClauses, " OR "), args...)
	}

	err := query.Find(&entities).Error
	return entities, total, err
}

func (r *BaseRepository[T]) FindByName(name string) (*T, error) {
	var entity T
	err := r.db.Where("name = ?", name).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindOneNotID(field, value string, excludeID string) (*T, error) {
	var entity T
	err := r.db.
		Where(fmt.Sprintf("%s = ?", field), value).
		Where("id != ?", excludeID).
		First(&entity).Error

	if err != nil {
		return nil, err
	}
	return &entity, nil
}
