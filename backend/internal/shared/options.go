package shared

import "gorm.io/gorm"

// FindOption interface untuk memodifikasi query GORM
type FindOption interface {
	Apply(*gorm.DB) *gorm.DB
}

// WithDeleted option untuk include soft deleted records
type WithDeleted struct{}

func (o WithDeleted) Apply(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}

// WithPreload option untuk preload relations
type WithPreload struct {
	Relations []string
}

func (o WithPreload) Apply(db *gorm.DB) *gorm.DB {
	for _, relation := range o.Relations {
		db = db.Preload(relation)
	}
	return db
}

// WithOrder option untuk mengurutkan hasil
type WithOrder struct {
	Order string
}

func (o WithOrder) Apply(db *gorm.DB) *gorm.DB {
	return db.Order(o.Order)
}

// WithLimit option untuk membatasi hasil
type WithLimit struct {
	Limit int
}

func (o WithLimit) Apply(db *gorm.DB) *gorm.DB {
	return db.Limit(o.Limit)
}

// WithOffset option untuk offset hasil
type WithOffset struct {
	Offset int
}

func (o WithOffset) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(o.Offset)
}

// WithWhere option untuk where condition
type WithWhere struct {
	Query interface{}
	Args  []interface{}
}

func (o WithWhere) Apply(db *gorm.DB) *gorm.DB {
	return db.Where(o.Query, o.Args...)
}

// WithJoin option untuk join tables
type WithJoin struct {
	Query string
	Args  []interface{}
}

func (o WithJoin) Apply(db *gorm.DB) *gorm.DB {
	return db.Joins(o.Query, o.Args...)
}

// WithSelect option untuk select specific fields
type WithSelect struct {
	Fields string
}

func (o WithSelect) Apply(db *gorm.DB) *gorm.DB {
	return db.Select(o.Fields)
}

// WithGroup option untuk group by
type WithGroup struct {
	Group string
}

func (o WithGroup) Apply(db *gorm.DB) *gorm.DB {
	return db.Group(o.Group)
}

// WithHaving option untuk having clause
type WithHaving struct {
	Query interface{}
	Args  []interface{}
}

func (o WithHaving) Apply(db *gorm.DB) *gorm.DB {
	return db.Having(o.Query, o.Args...)
}
