package domain

import "time"

type RefreshToken struct {
	BaseModel
	UserId    string    `gorm:"type:uuid"`
	Token     string    `gorm:"type:text"`
	ExpiresAt time.Time `gorm:"not null"`
}
