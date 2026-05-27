package domain

import "time"

type BaseModel struct {
	Id        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	CreatedBy string     `gorm:"not null"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	UpdatedBy *string
}
