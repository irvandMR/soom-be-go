package usecase

import (
	"errors"
	"soom-be-go/internal/constants"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/middleware"
	"soom-be-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type IngredientUsecase struct {
	repoIngredient repository.IngredientRepository
	repoCategory   repository.CategoriesRepository
	repoUom        repository.UomRepository
	repoHistory    repository.IngredientStockHistoryRepository
}

func NewIngredientUsecase(repoIngredient repository.IngredientRepository, repoCategory repository.CategoriesRepository, repoUom repository.UomRepository, repoHistory repository.IngredientStockHistoryRepository) *IngredientUsecase {
	return &IngredientUsecase{
		repoIngredient: repoIngredient,
		repoCategory:   repoCategory,
		repoUom:        repoUom,
		repoHistory:    repoHistory,
	}
}

func (u *IngredientUsecase) GetAllIngredient(req domain.PaginationRequest) (*domain.PaginationResponse, error) {
	req.Normalize()

	ingredient, total, err := u.repoIngredient.FindDataWithPagination(req)
	if err != nil {
		return nil, err
	}

	totalPage := int(total) / req.Limit

	if int(total)%req.Limit != 0 {
		totalPage++
	}

	return &domain.PaginationResponse{
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPage,
		Data:       ingredient,
	}, nil
}

func (u *IngredientUsecase) GetIngredientById(id string) (*domain.Ingredient, error) {
	ingredients, err := u.repoIngredient.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	return ingredients, nil
}

func (u *IngredientUsecase) CreateIngredient(req domain.IngredientRequest) (*domain.Ingredient, error) {

	uom, err := u.validateUom(req.UnitId)

	category, err := u.validateCategory(req.CategoryId)

	ingredient := &domain.Ingredient{
		BaseModelWithDeleted: domain.BaseModelWithDeleted{
			BaseModel: domain.BaseModel{
				CreatedAt: time.Now(),
				CreatedBy: req.Username,
			},
		},
		TenantId:     req.TenantId,
		CategoryId:   category.Id,
		UnitId:       uom.Id,
		Name:         req.Name,
		MinimumStock: req.MinimumStock,
		IsActive:     req.IsActive,
	}
	err = u.repoIngredient.Create(ingredient)
	return ingredient, err
}

func (u *IngredientUsecase) UpdateIngredient(req domain.IngredientRequestUpdate) (*domain.Ingredient, error) {
	ingredient, err := u.repoIngredient.FindById(req.Id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	ingredient.CategoryId = req.CategoryId
	ingredient.UnitId = req.UnitId
	ingredient.MinimumStock = req.MinimumStock
	ingredient.IsActive = req.IsActive
	ingredient.UpdatedAt = &now
	ingredient.UpdatedBy = &req.Username

	err = u.repoIngredient.Update(ingredient)
	return ingredient, err
}

func (u *IngredientUsecase) DeleteIngredient(id string, deletedby string) error {
	_, err := u.repoIngredient.FindById(id)

	if err != nil {
		return err
	}

	return u.repoIngredient.Delete(id, deletedby)
}

func (u *IngredientUsecase) StockIn(req domain.StockInRequest) (*domain.Ingredient, error) {
	ingredient, err := u.GetIngredientById(req.IngredientId)
	if err != nil {
		return nil, err
	}

	history := &domain.IngredientsStockHistory{
		BaseModel: domain.BaseModel{
			CreatedAt: time.Now(),
			CreatedBy: req.Username,
		},
		IngredientId:  ingredient.Id,
		Type:          string(constants.StockIn),
		Quantity:      req.Quantity,
		PurchasePrice: req.PurchasePrice,
		Notes:         req.Notes,
	}

	if err := u.repoHistory.Create(history); err != nil {
		return nil, err
	}

	// Disini insert cashflow

	oldStock := *ingredient.StockQuantity
	newStock := oldStock + req.Quantity

	var newAvgPrice float64
	if oldStock == 0 {
		newAvgPrice = req.PurchasePrice
	} else {
		oldStock := oldStock * *ingredient.AveragePrice
		newTotal := req.Quantity * req.PurchasePrice
		newAvgPrice = (oldStock + newTotal) / newStock
	}

	ingredient.StockQuantity = &newStock
	ingredient.AveragePrice = &newAvgPrice
	ingredient.PurchasePrice = &req.PurchasePrice

	if err := u.repoIngredient.Update(ingredient); err != nil {
		return nil, err
	}

	return ingredient, nil
}

func (u *IngredientUsecase) GetHistory(ingredientId string) ([]domain.IngredientsStockHistoryResponse, error) {
	_, err := u.repoIngredient.FindById(ingredientId)
	if err != nil {
		return nil, err
	}

	histories, err := u.repoHistory.FindByIngredientId(ingredientId)
	if err != nil {
		return nil, err
	}

	var responses []domain.IngredientsStockHistoryResponse
	for _, history := range histories {
		responses = append(responses, domain.IngredientsStockHistoryResponse{
			Quantity:      history.Quantity,
			PurchasePrice: history.PurchasePrice,
			Notes:         history.Notes,
			CreatedAt:     history.CreatedAt,
		})
	}

	return responses, nil
}

func (u *IngredientUsecase) validateUom(unitId string) (*domain.Uom, error) {
	uom, err := u.repoUom.FindById(unitId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &middleware.GlobalError{
				Message: "Unit not found",
				Code:    middleware.ErrNotFound.Code,
				Status:  middleware.ErrNotFound.Status,
			}
		}
		return nil, err
	}

	return uom, nil
}

func (u *IngredientUsecase) validateCategory(unitId string) (*domain.Categories, error) {
	category, err := u.repoCategory.FindById(unitId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &middleware.GlobalError{
				Message: "Category not found",
				Code:    middleware.ErrNotFound.Code,
				Status:  middleware.ErrNotFound.Status,
			}
		}
		return nil, err
	}

	return category, nil
}
