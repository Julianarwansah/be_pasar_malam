package models

import "time"

type Order struct {
	ID              uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint        `gorm:"index;not null" json:"user_id"`
	TotalAmount     float64     `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status          string      `gorm:"size:32;default:'pending'" json:"status"`
	ShippingAddress string      `gorm:"type:text" json:"shipping_address"`
	Notes           string      `gorm:"type:text" json:"notes"`
	PaymentMethod   string      `gorm:"size:32;not null" json:"payment_method"`
	VANumber        string      `gorm:"size:64" json:"va_number,omitempty"`
	GopayDeeplink   string      `gorm:"size:512" json:"gopay_deeplink,omitempty"`
	Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID     uint    `gorm:"index;not null" json:"order_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `gorm:"size:191" json:"product_name"`
	Price       float64 `gorm:"type:decimal(12,2)" json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `gorm:"type:decimal(12,2)" json:"subtotal"`
}
