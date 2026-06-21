package repository

import (
	"soom-be-go/internal/domain"

	"gorm.io/gorm"
)

type ProductRecipeItemRepository interface {
	FindByRecipeId(recipeId string) ([]domain.ProductRecipeItem, error)
}

type productRecipeItemRepository struct {
	db *gorm.DB
}

// FindByRecipeId implements [ProductRecipeItemRepository].
func (p *productRecipeItemRepository) FindByRecipeId(recipeId string) ([]domain.ProductRecipeItem, error) {
	var items []domain.ProductRecipeItem
	result := p.db.Where("recipe_id = ?", recipeId).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func NewProductRecipeItemRepository(db *gorm.DB) ProductRecipeItemRepository {
	return &productRecipeItemRepository{db: db}
}
