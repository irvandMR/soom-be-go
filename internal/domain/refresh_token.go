package domain

import "time"

type RefreshToken struct {
	BaseModel
	UserId    string    `gorm:"type:uuid"`
	Token     string    `gorm:"size:255;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
