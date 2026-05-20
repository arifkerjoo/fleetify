package repositories

import "backend/internal/shared"

// BaseRepositoryInterface adalah kontrak umum untuk semua repository
// menggunakan generic agar bisa dipakai ulang untuk semua entitas.
type BaseRepositoryInterface[T any] interface {
	FindAllPaginated(limit, offset int, search string, field []string) ([]T, int64, error)
	FindAll(conditions map[string]interface{}, options ...shared.FindOption) ([]T, error)
	FindOne(conditions map[string]interface{}, options ...shared.FindOption) (*T, error)
	FindByName(name string) (*T, error)
	FindOneNotID(fields, value string, excludeID string) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id string) error
}
