package repository

import (
	"soom-be-go/internal/domain"

	"gorm.io/gorm"
)

type IngredientStockHistoryRepository interface {
	Create(req *domain.IngredientsStockHistory) error
	FindByIngredientId(id string) ([]domain.IngredientsStockHistory, error)
}

type ingredientStockRepository struct {
	db *gorm.DB
}

// FindByIngredientId implements [IngredientStockHistoryRepository].
func (i *ingredientStockRepository) FindByIngredientId(id string) ([]domain.IngredientsStockHistory, error) {
	var histories []domain.IngredientsStockHistory
	result := i.db.Where("ingredient_id = ?", id).Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	return histories, nil

}

// Create implements [IngredientStockHistoryRepository].
func (i *ingredientStockRepository) Create(req *domain.IngredientsStockHistory) error {
	return i.db.Create(req).Error
}

func NewIngredientStockRepository(db *gorm.DB) IngredientStockHistoryRepository {
	return &ingredientStockRepository{db: db}
}
