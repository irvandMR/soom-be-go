package repository

import (
	"soom-be-go/internal/domain"
	"time"

	"gorm.io/gorm"
)

type CategoriesRepository interface {
	FindAll() ([]domain.Categories, error)
	FindDataWithPagination(req domain.PaginationRequest) ([]domain.Categories, int64, error)
	Create(req *domain.Categories) error
	Update(req *domain.Categories) error
	Delete(id string, deletedBy string) error
	FindById(id string) (*domain.Categories, error)
}

type categoriesRepository struct {
	db *gorm.DB
}

func (r *categoriesRepository) FindAll() ([]domain.Categories, error) {
	var categories []domain.Categories
	result := r.db.Where("deleted_at is null").Find(&categories)
	return categories, result.Error
}

func (r *categoriesRepository) FindDataWithPagination(req domain.PaginationRequest) ([]domain.Categories, int64, error) {
	var categories []domain.Categories
	var total int64

	query := r.db.Model(&domain.Categories{}).Where("deleted_at is null")

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
