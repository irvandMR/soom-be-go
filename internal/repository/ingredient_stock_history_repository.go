package repository

import (
	"soom-be-go/internal/domain"

	"gorm.io/gorm"
)

type IngredientStockHistoryRepository interface {
	Create(req *domain.IngredientsStockHistory) error
	FindByIngredientId(id string) ([]domain.IngredientsStockHistory, error)
	FindAllHistory(req domain.IngredientsStockHistoryRequest) ([]domain.IngredientsStockHistory, error)
	FindById(id string) (*domain.IngredientsStockHistory, error)
	Update(req *domain.IngredientsStockHistory) error
}

type ingredientStockRepository struct {
	db *gorm.DB
}

// FindAllHistory implements [IngredientStockHistoryRepository].
func (i *ingredientStockRepository) FindAllHistory(req domain.IngredientsStockHistoryRequest) ([]domain.IngredientsStockHistory, error) {
	var history []domain.IngredientsStockHistory

	query := i.db.Model(&domain.IngredientsStockHistory{}).Where("ingredient_id = ?", req.IngredientId)

	if req.StartDate != nil && req.EndDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	} else if req.StartDate != nil {
		query = query.Where("created_at >= ?", req.StartDate)
	} else if req.EndDate != nil {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	if err := query.Find(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
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

// FindById implements [IngredientStockHistoryRepository].
func (i *ingredientStockRepository) FindById(id string) (*domain.IngredientsStockHistory, error) {
	var history domain.IngredientsStockHistory
	result := i.db.Where("id = ?", id).First(&history)
	if result.Error != nil {
		return nil, result.Error
	}
	return &history, nil
}

// Update implements [IngredientStockHistoryRepository].
func (i *ingredientStockRepository) Update(req *domain.IngredientsStockHistory) error {
	return i.db.Save(req).Error
}
