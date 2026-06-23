package repository

import (
	"soom-be-go/internal/domain"

	"gorm.io/gorm"
)

type ProductRecipeItemRepository interface {
	FindByRecipeId(recipeId string) ([]domain.ProductRecipeItem, error)
	CreateWithItems(recipe *domain.ProductRecipes, items []domain.ProductRecipeItem) error
}

type productRecipeItemRepository struct {
	db *gorm.DB
}

// CreateWithItems implements [ProductRecipeItemRepository].
func (p *productRecipeItemRepository) CreateWithItems(recipe *domain.ProductRecipes, items []domain.ProductRecipeItem) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create recipe
		if err := tx.Create(recipe).Error; err != nil {
			return err
		}

		// 2. Insert all items
		for i := range items {
			items[i].RecipeId = recipe.Id
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
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
