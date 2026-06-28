package models

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CartID    uint      `gorm:"index;not null" json:"cart_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Subtotal  float64   `gorm:"type:decimal(12,2)" json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
