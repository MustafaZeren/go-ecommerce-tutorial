// Package models veritabanı şemalarını ve uygulama modellerini içerir.
package models

import (
	"time"

	"gorm.io/gorm"
)

// Order müşterilerin verdiği siparişlerin ana bilgilerini tutar.
type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Status      OrderStatus    `json:"status" gorm:"default:pending"`
	TotalAmount float64        `json:"total_amount" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User       User        `json:"user"`
	OrderItems []OrderItem `json:"order_items"`
}

// OrderStatus siparişin mevcut durumunu belirten özel bir tiptir.
type OrderStatus string

const (
	// OrderStatusPending siparişin onay beklediğini ifade eder.
	OrderStatusPending OrderStatus = "pending"
	// OrderStatusConfirmed siparişin ödemesinin alındığını ve onaylandığını ifade eder.
	OrderStatusConfirmed OrderStatus = "confirmed"
	// OrderStatusShipped siparişin kargoya verildiğini ifade eder.
	OrderStatusShipped OrderStatus = "shipped"
	// OrderStatusDelivered siparişin müşteriye teslim edildiğini ifade eder.
	OrderStatusDelivered OrderStatus = "delivered"
	// OrderStatusCancelled siparişin iptal edildiğini ifade eder.
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem bir sipariş içerisindeki her bir ürünün miktar ve fiyat bilgisini tutar.
type OrderItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OrderID   uint           `json:"order_id" gorm:"not null"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	Price     float64        `json:"price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Order   Order   `json:"-"`
	Product Product `json:"product"`
}

// Cart kullanıcının aktif sepet bilgisini temsil eder.
type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	CartItems []CartItem `json:"cart_items"`
}

// CartItem sepete eklenen her bir ürünü ve miktarını temsil eder.
type CartItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CartID    uint           `json:"cart_id" gorm:"not null"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Cart    Cart    `json:"-"`
	Product Product `json:"product"`
}
