package repository

import (
	"soom-be-go/internal/domain"

	"gorm.io/gorm"
)

type TenantRepository interface {
	FindCodeTenantById(id string) (string, error)
}

type tenantRepository struct {
	db *gorm.DB
}

// FindCodeTenantById implements [TenantRepository].
func (t *tenantRepository) FindCodeTenantById(id string) (string, error) {
	var code string
	err := t.db.Model(&domain.Tenant{}).Select("code").Where("id = ?", id).First(&code).Error

	return code, err
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{
		db: db,
	}
}
