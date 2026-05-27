package usecase

import (
	"errors"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/middleware"
	"soom-be-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type IngredientUsecase struct {
	repo repository.IngredientRepository
}

func NewIngredientUsecase(r repository.IngredientRepository) *IngredientUsecase {
	return &IngredientUsecase{repo: r}
}

func (u *IngredientUsecase) GetAllIngredient(req domain.PaginationRequest) (*domain.PaginationResponse, error) {
	req.Normalize()

	ingredient, total, err := u.repo.FindDataWithPagination(req)
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
	ingredients, err := u.repo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	return ingredients, nil
}

func (u *IngredientUsecase) CreateIngredient(req domain.IngredientRequest) (*domain.Ingredient, error) {

	ingredient := &domain.Ingredient{
		BaseModelWithDeleted: domain.BaseModelWithDeleted{
			BaseModel: domain.BaseModel{
				CreatedAt: time.Now(),
				CreatedBy: req.Username,
			},
		},
		TenantId:      req.TenantId,
		CategoryId:    req.CategoryId,
		UnitId:        req.UnitId,
		Name:          req.Name,
		StockQuantity: *req.StockQuantity,
		MinimumStock:  req.MinimumStock,
		PricePerUnit:  req.PricePerUnit,
		PurchasePrice: req.PurchasePrice,
		AveragePrice:  req.AveragePrice,
		IsActive:      req.IsActive,
	}
	err := u.repo.Create(*ingredient)
	return ingredient, err
}

func (u *IngredientUsecase) UpdateIngredient(req domain.IngredientRequestUpdate) (*domain.Ingredient, error) {
	ingredient, err := u.repo.FindById(req.Id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	ingredient.CategoryId = req.CategoryId
	ingredient.UnitId = req.UnitId
	ingredient.MinimumStock = req.MinimumStock
	ingredient.PurchasePrice = req.PurchasePrice
	ingredient.IsActive = req.IsActive
	ingredient.UpdatedAt = &now
	ingredient.UpdatedBy = &req.Username

	err = u.repo.Update(ingredient)
	return ingredient, err
}

func (u *IngredientUsecase) DeleteIngredient(id string, deletedby string) error {
	_, err := u.repo.FindById(id)

	if err != nil {
		return err
	}

	return u.repo.Delete(id, deletedby)
}
