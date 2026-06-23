package repository

import (
	"soom-be-go/internal/domain"

	"gorm.io/gorm"
)

type ProductRecipeRepository interface {
	FindByProductId(productId string) ([]domain.ProductRecipes, error)
	Update(recipe *domain.ProductRecipes) error
	FindByProductIdAndIsActiveTrue(productId string) (*domain.ProductRecipes, error)
	CountByProductId(productId string) (int64, error)
}

type productRecipeRepository struct {
	db *gorm.DB
}

// CountByProductId implements [ProductRecipeRepository].
func (p *productRecipeRepository) CountByProductId(productId string) (int64, error) {
	var count int64
	err := p.db.Model(&domain.ProductRecipes{}).Where("product_id = ?", productId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindByProductIdAndIsActiveTrue implements [ProductRecipeRepository].
func (p *productRecipeRepository) FindByProductIdAndIsActiveTrue(productId string) (*domain.ProductRecipes, error) {
	var recipe domain.ProductRecipes
	result := p.db.Where("product_id = ? AND is_active = true", productId).First(&recipe)
	if result.Error != nil {
		return nil, result.Error
	}
	return &recipe, nil
}

// Update implements [ProductRecipeRepository].
func (p *productRecipeRepository) Update(recipe *domain.ProductRecipes) error {
	result := p.db.Save(recipe)
	return result.Error
}

// FindByProductId implements [ProductRecipeRepository].
func (p *productRecipeRepository) FindByProductId(productId string) ([]domain.ProductRecipes, error) {
	var recipe []domain.ProductRecipes
	result := p.db.Where("product_id = ?", productId).Find(&recipe)
	if result.Error != nil {
		return nil, result.Error
	}
	return recipe, nil
}

func NewProductRecipeRepository(db *gorm.DB) ProductRecipeRepository {
	return &productRecipeRepository{db: db}
}
