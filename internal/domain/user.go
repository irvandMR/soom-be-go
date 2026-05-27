package domain

import "time"

type User struct {
	BaseModelWithDeleted
	TenantId              *string    `gorm:"type:uuid"`
	TenantRole            string     `gorm:"size:50;not null"`
	Username              string     `gorm:"size:50;not null;unique"`
	Email                 string     `gorm:"not null;unique"`
	Password              string     `gorm:"not null"`
	Role                  string     `gorm:"size:50;not null"`
	MustChangePassword    bool       `gorm:"default:true;not null"`
	TempPasswordExpiresAt *time.Time `gorm:"type:timestamptz"`
	IsActive              bool       `gorm:"default:true;not null"`
}
