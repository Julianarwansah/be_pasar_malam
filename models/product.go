package models

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:191;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	Stock       int       `gorm:"default:0" json:"stock"`
	Category    string    `gorm:"size:64;index" json:"category"`
	ImageURL    string    `gorm:"size:512" json:"image_url"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
