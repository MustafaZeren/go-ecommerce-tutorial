// Package models veritabanı şemalarını ve uygulama modellerini içerir.
package models

import (
	"time"

	"gorm.io/gorm"
)

// User sistemdeki kullanıcı hesap bilgilerini temsil eder.
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Phone     string         `json:"phone"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	Role      UserRole       `json:"role" gorm:"default:customer"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	RefreshTokens []RefreshToken `json:"-"`
	Orders        []Order        `json:"-"`
	Cart          Cart           `json:"-"`
}

// UserRole kullanıcı yetkilendirme seviyelerini belirleyen özel tiptir.
type UserRole string

const (
	// UserRoleCustomer standart alışveriş yapan kullanıcı rolüdür.
	UserRoleCustomer UserRole = "customer"
	// UserRoleAdmin tüm yönetim yetkilerine sahip yönetici rolüdür.
	UserRoleAdmin UserRole = "admin"
)

// RefreshToken JWT oturumlarının yenilenmesi için kullanılan belirteç bilgisini tutar.
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"-"`
}
