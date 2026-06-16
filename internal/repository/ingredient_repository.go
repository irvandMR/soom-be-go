package repository

import (
	"soom-be-go/internal/domain"
	"time"

	"gorm.io/gorm"
)

type IngredientRepository interface {
	FindAll(req domain.IngredientQueryRequest) ([]domain.Ingredient, int64, error)
	FindById(id string) (*domain.Ingredient, error)
	Create(req *domain.Ingredient) error
	Update(req *domain.Ingredient) error
	Delete(id string, deletedBy string) error
	FindAllNoPagination(tenantId *string) ([]domain.Ingredient, error)
}

type ingredientRepository struct {
	db *gorm.DB
}

// FindAllNoPagination implements [IngredientRepository].
func (r *ingredientRepository) FindAllNoPagination(tenantId *string) ([]domain.Ingredient, error) {
	var ingredient []domain.Ingredient
	result := r.db.Preload("Category").Preload("Unit").Where("tenant_id = ? and deleted_at is null", tenantId).Find(&ingredient)
	if result.Error != nil {
		return nil, result.Error
	}
	return ingredient, nil
}

// FindAll implements [IngredientRepository].
func (r *ingredientRepository) FindAll(req domain.IngredientQueryRequest) ([]domain.Ingredient, int64, error) {
	var ingredient []domain.Ingredient
	var total int64

	query := r.db.Model(&domain.Ingredient{}).Where("deleted_at is null")

	if req.TenantId != nil {
		query = query.Where("tenant_id = ?", *req.TenantId)
	}

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if req.CategoriesId != "" {
		query = query.Where("category_id = ?", req.CategoriesId)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)

	result := query.Preload("Category").Preload("Unit").Limit(req.Limit).Offset(req.Offset()).Find(&ingredient)
	return ingredient, total, result.Error
}

// Delete implements [IngredientRepository].
func (r *ingredientRepository) Delete(id string, deletedBy string) error {
	now := time.Now()
	result := r.db.Model(&domain.Ingredient{}).Where("id = ? and deleted_at is null", id).Updates(map[string]interface{}{
		"deleted_by": deletedBy,
		"deleted_at": now,
	})
	return result.Error
}

func (r *ingredientRepository) FindById(id string) (*domain.Ingredient, error) {
	var ingredient domain.Ingredient
	result := r.db.Preload("Category").Preload("Unit").Where("id = ? and deleted_at is null", id).First(&ingredient)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ingredient, nil
}

func (r *ingredientRepository) Create(req *domain.Ingredient) error {
	return r.db.Create(req).Error
}

func (r *ingredientRepository) Update(req *domain.Ingredient) error {
	return r.db.Save(req).Error
}
func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return &ingredientRepository{db: db}
}
