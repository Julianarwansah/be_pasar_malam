package models

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FirebaseUID   string    `gorm:"size:191;uniqueIndex;not null" json:"firebase_uid"`
	Email         string    `gorm:"size:191;uniqueIndex;not null" json:"email"`
	Name          string    `gorm:"size:191" json:"name"`
	Role          string    `gorm:"size:32;default:'customer'" json:"role"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	FCMToken      string    `gorm:"type:text" json:"fcm_token,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
