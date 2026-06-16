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
	db             *gorm.DB
	repoIngredient repository.IngredientRepository
	repoCategory   repository.CategoriesRepository
	repoUom        repository.UomRepository
	repoHistory    repository.IngredientStockHistoryRepository
}

func NewIngredientUsecase(db *gorm.DB, repoIngredient repository.IngredientRepository, repoCategory repository.CategoriesRepository, repoUom repository.UomRepository, repoHistory repository.IngredientStockHistoryRepository) *IngredientUsecase {
	return &IngredientUsecase{
		db:             db,
		repoIngredient: repoIngredient,
		repoCategory:   repoCategory,
		repoUom:        repoUom,
		repoHistory:    repoHistory,
	}
}

func (u *IngredientUsecase) GetAllIngredient(req domain.IngredientQueryRequest) (*domain.PaginationResponse, error) {
	req.Normalize()

	ingredient, total, err := u.repoIngredient.FindAll(req)
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
		Data:       u.mappingResponse(ingredient),
	}, nil
}

func (u *IngredientUsecase) GetAllIngredientWithoutPaginaton(tenantId *string) ([]domain.IngredientResponse, error) {
	ingredients, err := u.repoIngredient.FindAllNoPagination(tenantId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	result := u.mappingResponse(ingredients)
	return result, nil
}

func (u *IngredientUsecase) GetIngredientById(id string) (*domain.IngredientResponse, error) {
	ingredients, err := u.repoIngredient.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	return &u.mappingResponse([]domain.Ingredient{*ingredients})[0], nil
}

func (u *IngredientUsecase) CreateIngredient(req domain.IngredientRequest) (*domain.Ingredient, error) {

	uom, err := u.validateUom(req.UnitId)
	if err != nil {
		return nil, err
	}

	category, err := u.validateCategory(req.CategoryId)
	if err != nil {
		return nil, err
	}

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

	_, err = u.validateUom(req.UnitId)
	if err != nil {
		return nil, err
	}

	_, err = u.validateCategory(req.CategoryId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	ingredient.CategoryId = req.CategoryId
	ingredient.UnitId = req.UnitId
	ingredient.Name = req.Name
	ingredient.MinimumStock = req.MinimumStock

	if req.HistoryId != "" {
		history, err := u.repoHistory.FindById(req.HistoryId)
		if err != nil {
			return nil, err
		}

		newAvgPrice := 0.0
		if ingredient.AveragePrice != nil {
			newAvgPrice = *ingredient.AveragePrice
		}

		if ingredient.StockQuantity != nil && *ingredient.StockQuantity > 0 {
			oldTotalValue := *ingredient.StockQuantity * newAvgPrice
			priceDifference := req.PurchasePrice - history.PurchasePrice
			newTotalValue := oldTotalValue + (priceDifference * history.Quantity)
			if newTotalValue < 0 {
				newTotalValue = 0
			}
			newAvgPrice = newTotalValue / *ingredient.StockQuantity
		}

		ingredient.AveragePrice = &newAvgPrice

		history.PurchasePrice = req.PurchasePrice
		if err := u.repoHistory.Update(history); err != nil {
			return nil, err
		}
	}

	ingredient.IsActive = req.IsActive
	ingredient.UpdatedAt = &now
	ingredient.UpdatedBy = &req.Username

	// Recheck status when MinimumStock changes
	if ingredient.StockQuantity != nil {
		if ingredient.MinimumStock >= *ingredient.StockQuantity {
			ingredient.Status = string(constants.Critical)
		} else {
			ingredient.Status = string(constants.Save)
		}
	}

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
	ingredient, err := u.repoIngredient.FindById(req.IngredientId)
	if err != nil {
		return nil, err
	}

	// safe ambil stock lama
	oldStock := 0.0
	if ingredient.StockQuantity != nil {
		oldStock = *ingredient.StockQuantity
	}

	newStock := oldStock + req.Quantity

	// safe ambil average price lama
	oldAvg := 0.0
	if ingredient.AveragePrice != nil {
		oldAvg = *ingredient.AveragePrice
	}

	var newAvgPrice float64

	if oldStock == 0 {
		newAvgPrice = req.PurchasePrice
	} else {
		oldTotal := oldStock * oldAvg
		newTotal := req.Quantity * req.PurchasePrice
		newAvgPrice = (oldTotal + newTotal) / newStock
	}

	ingredient.StockQuantity = &newStock
	ingredient.AveragePrice = &newAvgPrice
	ingredient.PurchasePrice = &req.PurchasePrice
	ingredient.Status = string(constants.Save)

	// jika stock <= minimum stock maka status menjadi critical
	if ingredient.MinimumStock >= newStock {
		ingredient.Status = string(constants.Critical)
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

	// TODO: Disini insert cashflow

	if err := u.repoIngredient.Update(ingredient); err != nil {
		return nil, err
	}

	return ingredient, nil
}

func (u *IngredientUsecase) GetHistory(req domain.IngredientsStockHistoryRequest) ([]domain.IngredientsStockHistoryResponse, error) {

	if err := u.validationRangeMonth(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	_, err := u.repoIngredient.FindById(req.IngredientId)
	if err != nil {
		return nil, err
	}

	histories, err := u.repoHistory.FindAllHistory(req)
	if err != nil {
		return nil, err
	}

	var responses []domain.IngredientsStockHistoryResponse
	for _, history := range histories {
		responses = append(responses, domain.IngredientsStockHistoryResponse{
			Type:          history.Type,
			Quantity:      history.Quantity,
			PurchasePrice: history.PurchasePrice,
			Notes:         history.Notes,
			CreatedAt:     history.CreatedAt,
			Id:            history.Id,
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

func (u *IngredientUsecase) validateCategory(categoryId string) (*domain.Categories, error) {
	category, err := u.repoCategory.FindById(categoryId)
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

func (u *IngredientUsecase) mappingResponse(ingredients []domain.Ingredient) []domain.IngredientResponse {
	result := make([]domain.IngredientResponse, 0, len(ingredients))

	for _, ingredient := range ingredients {
		result = append(result, domain.IngredientResponse{
			Id:            ingredient.Id,
			CategoryId:    ingredient.CategoryId,
			UnitId:        ingredient.UnitId,
			Name:          ingredient.Name,
			StockQuantity: ingredient.StockQuantity,
			MinimumStock:  ingredient.MinimumStock,
			PricePerUnit:  ingredient.PricePerUnit,
			PurchasePrice: ingredient.PurchasePrice,
			AveragePrice:  ingredient.AveragePrice,
			IsActive:      ingredient.IsActive,
			Status:        string(ingredient.Status),
			CategoryName:  ingredient.Category.Name,
			UnitSymbol:    ingredient.Unit.Symbol,
		})
	}
	return result
}

func (u *IngredientUsecase) validationRangeMonth(startDate *time.Time, endDate *time.Time) error {
	if startDate != nil && endDate != nil {
		diff := endDate.Sub(*startDate)
		if diff > 31*24*time.Hour {
			return &middleware.GlobalError{
				Message: "Date range cannot exceed 1 month",
				Code:    middleware.ErrUnprocessableEntity.Code,
				Status:  middleware.ErrUnprocessableEntity.Status,
			}
		}
	}
	return nil
}
