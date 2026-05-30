package repository

import (
	"soom-be-go/internal/domain"
	"time"

	"gorm.io/gorm"
)

type IngredientRepository interface {
	FindAll() ([]domain.Ingredient, error)
	FindById(id string) (*domain.Ingredient, error)
	Create(req *domain.Ingredient) error
	Update(req *domain.Ingredient) error
	Delete(id string, deletedBy string) error
	FindDataWithPagination(req domain.PaginationRequest) ([]domain.Ingredient, int64, error)
}

type ingredientRepository struct {
	db *gorm.DB
}

// FindDataWithPagination implements [IngredientRepository].
func (r *ingredientRepository) FindDataWithPagination(req domain.PaginationRequest) ([]domain.Ingredient, int64, error) {
	var ingredient []domain.Ingredient
	var total int64

	query := r.db.Model(&domain.Ingredient{}).Where("deleted_at is null")

	query.Count(&total)

	result := query.Limit(req.Limit).Offset(req.Offset()).Find(&ingredient)
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

func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return &ingredientRepository{db: db}
}

func (r *ingredientRepository) FindAll() ([]domain.Ingredient, error) {
	var ingredient []domain.Ingredient
	result := r.db.Where("deleted_at is null").Find(&ingredient)
	return ingredient, result.Error
}

func (r *ingredientRepository) FindById(id string) (*domain.Ingredient, error) {
	var ingredient domain.Ingredient
	result := r.db.Where("id = ? and deleted_at is null", id).First(&ingredient)
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
