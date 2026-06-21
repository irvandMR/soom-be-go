package repository

import (
	"soom-be-go/internal/domain"

	"gorm.io/gorm"
)

type ProductRecipeRepository interface {
	FindByProductId(productId string) ([]domain.ProductRecipes, error)
}

type productRecipeRepository struct {
	db *gorm.DB
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
