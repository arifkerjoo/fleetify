package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "SA"
	RoleApproval   UserRole = "APPROVAL"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type User struct {
	BaseModel

	Username   string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_username" json:"username"`
	Email      string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_tenant_email" json:"email"`
	Password   string     `gorm:"type:varchar(255);not null" json:"-"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	FullName string   `gorm:"type:varchar(255);not null" json:"full_name"`
	Phone    string   `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Role     UserRole `gorm:"type:enum('SA','APPROVAL');index;not null" json:"role"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}

type UserResponse struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	Username   string     `json:"username"`
	FullName   string     `json:"full_name"`
	Phone      string     `json:"phone,omitempty"`
	Role       UserRole   `json:"role"`
	IsActive   bool       `json:"is_active"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type UserListResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Username:   u.Username,
		FullName:   u.FullName,
		Phone:      u.Phone,
		Role:       u.Role,
		IsActive:   u.IsActive,
		VerifiedAt: u.VerifiedAt,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (u *User) ToListResponse() *UserListResponse {
	return &UserListResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
