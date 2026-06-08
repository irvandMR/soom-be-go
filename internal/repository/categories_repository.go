package repository

import (
	"soom-be-go/internal/domain"
	"time"

	"gorm.io/gorm"
)

type CategoriesRepository interface {
	FindAll(req domain.CategoriesQueryRequest) ([]domain.Categories, int64, error)
	Create(req *domain.Categories) error
	Update(req *domain.Categories) error
	Delete(id string, deletedBy string) error
	FindById(id string) (*domain.Categories, error)
	FindTypeByTenant(tenantId *string) ([]domain.Categories, error)
}

type categoriesRepository struct {
	db *gorm.DB
}

// FindTypeByTenant implements [CategoriesRepository].
func (r *categoriesRepository) FindTypeByTenant(tenantId *string) ([]domain.Categories, error) {
	var categories []domain.Categories

	query := r.db.Model(&domain.Categories{})

	if tenantId != nil {
		query = query.Where("tenant_id = ?", *tenantId)
	}

	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// FindAll implements [CategoriesRepository].
func (r *categoriesRepository) FindAll(req domain.CategoriesQueryRequest) ([]domain.Categories, int64, error) {
	var categories []domain.Categories
	var total int64

	query := r.db.Model(&domain.Categories{}).Where("deleted_at is null")

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	query.Count(&total)

	result := query.Limit(req.Limit).Offset(req.Offset()).Find(&categories)
	return categories, total, result.Error
}

// Create implements [CategoriesRepository].
func (r *categoriesRepository) Create(req *domain.Categories) error {
	return r.db.Create(req).Error
}

// Delete implements [CategoriesRepository].
func (r *categoriesRepository) Delete(id string, deletedBy string) error {
	now := time.Now()
	result := r.db.Model(&domain.Categories{}).Where("id = ? and deleted_at is null", id).Updates(map[string]interface{}{
		"deleted_by": deletedBy,
		"deleted_at": now,
	})
	return result.Error
}

// FindById implements [CategoriesRepository].
func (r *categoriesRepository) FindById(id string) (*domain.Categories, error) {
	var categories domain.Categories
	result := r.db.Where("id = ? and deleted_at is null", id).First(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return &categories, nil
}

// Update implements [CategoriesRepository].
func (r *categoriesRepository) Update(req *domain.Categories) error {
	result := r.db.Save(req)
	return result.Error
}

func NewCategoriesRepository(db *gorm.DB) CategoriesRepository {
	return &categoriesRepository{db: db}
}
