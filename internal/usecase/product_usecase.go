package usecase

import (
	"errors"
	"fmt"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/middleware"
	"soom-be-go/internal/pkg/uom"
	"soom-be-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type ProductUsecase struct {
	repo           repository.ProductRepository
	repoCategory   repository.CategoriesRepository
	repoUom        repository.UomRepository
	repoTenant     repository.TenantRepository
	repoRecipe     repository.ProductRecipeRepository
	repoIngredient repository.IngredientRepository
	repoItem       repository.ProductRecipeItemRepository
}

type recipeItemWithDetail struct {
	Item       domain.ProductRecipeItem
	Ingredient domain.Ingredient
	Unit       domain.Uom
}

func NewProductUsecase(repoProduct repository.ProductRepository, repoCategory repository.CategoriesRepository, repoUom repository.UomRepository, repoTenant repository.TenantRepository, repoRecipe repository.ProductRecipeRepository, repoIngredient repository.IngredientRepository, repoItem repository.ProductRecipeItemRepository) *ProductUsecase {
	return &ProductUsecase{repo: repoProduct, repoCategory: repoCategory, repoUom: repoUom, repoTenant: repoTenant, repoRecipe: repoRecipe, repoIngredient: repoIngredient, repoItem: repoItem}
}

func (p *ProductUsecase) DeleteProduct(id string, deletedby string) error {
	_, err := p.repo.FindById(id)

	if err != nil {
		return err
	}

	return p.repo.Delete(id, deletedby)
}

func (p *ProductUsecase) UpdateProduct(req domain.ProductRequestUpdate) (*domain.ProductResponse, error) {
	product, err := p.repo.FindById(req.Id)
	if err != nil {
		return nil, err
	}

	uom, err := p.validateUom(req.UnitID)
	if err != nil {
		return nil, err
	}

	category, err := p.validationCategory(req.CategoryID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	product.CategoryID = category.Id
	product.UnitID = uom.Id
	product.Name = req.Name
	product.IsActive = req.IsActive
	product.Type = req.Type
	product.UpdatedAt = &now
	product.UpdatedBy = &req.Username

	err = p.repo.Update(product)

	response := p.mapppingSingleResponse(*product)
	return &response, err

}

func (p *ProductUsecase) GetProductById(id string) (*domain.ProductResponse, error) {
	product, err := p.repo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	return &p.mappingResponse([]domain.Product{*product})[0], nil
}

func (p *ProductUsecase) GetAllProducts(req domain.ProductQueryRequest) (*domain.PaginationResponse, error) {
	req.Normalize()

	product, total, err := p.repo.FindAllPagination(req)
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
		Data:       p.mappingResponse(product),
	}, nil
}

func (p *ProductUsecase) CreateProduct(req domain.ProductRequest) (*domain.ProductResponse, error) {
	uom, err := p.validateUom(req.UnitID)
	if err != nil {
		return nil, err
	}

	category, err := p.validationCategory(req.CategoryID)
	if err != nil {
		return nil, err
	}

	code, err := p.repoTenant.FindCodeTenantById(*req.TenantId)
	if err != nil {
		return nil, err
	}

	productCode := fmt.Sprintf("%s-%s", code, req.Code)

	product := &domain.Product{
		BaseModelWithDeleted: domain.BaseModelWithDeleted{
			BaseModel: domain.BaseModel{
				CreatedAt: time.Now(),
				CreatedBy: req.Username,
			},
		},
		CategoryID: category.Id,
		UnitID:     uom.Id,
		Code:       productCode,
		Name:       req.Name,
		Type:       req.Type,
		TenantID:   req.TenantId,
		IsActive:   req.IsActive,
	}

	err = p.repo.Create(product)

	product.Category = *category
	product.Unit = *uom

	response := p.mapppingSingleResponse(*product)
	return &response, err
}

func (p *ProductUsecase) SaveProductRecipe(req domain.ProductRecipesRequest) (*domain.RecipeResponse, error) {
	now := time.Now()

	product, err := p.repo.FindById(req.ProductId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &middleware.GlobalError{
				Message: "Product not found",
				Code:    middleware.ErrNotFound.Code,
				Status:  middleware.ErrNotFound.Status,
			}
		}
		return nil, err
	}

	// Validasi yield unit dari request (lihat pembahasan sebelumnya:
	// estimated_yield butuh satuannya sendiri, tidak diwarisi dari produk)
	yieldUnit, err := p.validateUom(req.UnitId)
	if err != nil {
		return nil, err
	}

	// Nonaktifkan resep aktif sebelumnya (kalau ada)
	activeRecipe, err := p.repoRecipe.FindByProductIdAndIsActiveTrue(product.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if activeRecipe != nil {
		activeRecipe.IsActive = false
		activeRecipe.UpdatedAt = &now
		username := req.Username
		activeRecipe.UpdatedBy = &username
		if err := p.repoRecipe.Update(activeRecipe); err != nil {
			return nil, err
		}
	}

	count, err := p.repoRecipe.CountByProductId(req.ProductId)
	if err != nil {
		return nil, err
	}
	newVersion := int32(count + 1)

	newRecipe := &domain.ProductRecipes{
		BaseModel: domain.BaseModel{
			CreatedAt: now,
			CreatedBy: req.Username,
		},
		ProductId:      req.ProductId,
		VersionNumber:  newVersion,
		Notes:          req.Notes,
		EstimatedYield: &req.EstimatedYield,
		UnitId:         &req.UnitId,
		IsActive:       true,
	}

	var totalCost float64
	items := make([]domain.ProductRecipeItem, 0, len(req.Items))
	itemsWithDetail := make([]recipeItemWithDetail, 0, len(req.Items))

	for _, itemReq := range req.Items {
		recipeUnit, err := p.validateUom(itemReq.UnitId)
		if err != nil {
			return nil, err
		}

		ingredient, err := p.validationIngredient(itemReq.IngredientId)
		if err != nil {
			return nil, err
		}

		stockUnit, err := p.validateUom(ingredient.UnitId)
		if err != nil {
			return nil, err
		}

		qtyInStockUnit, err := uom.Convert(itemReq.Quantity, *recipeUnit, *stockUnit)
		if err != nil {
			if errors.Is(err, uom.ErrIncompatibleUnit) {
				return nil, &middleware.GlobalError{
					Message: "Satuan resep tidak cocok dengan satuan stok bahan baku",
					Code:    middleware.ErrBadRequest.Code,
					Status:  middleware.ErrBadRequest.Status,
				}
			}
			return nil, err
		}

		if ingredient.AveragePrice == nil {
			return nil, &middleware.GlobalError{
				Message: "Bahan baku belum memiliki harga rata-rata",
				Code:    middleware.ErrBadRequest.Code,
				Status:  middleware.ErrBadRequest.Status,
			}
		}
		ingredientCost := qtyInStockUnit * *ingredient.AveragePrice
		totalCost += ingredientCost

		newItem := domain.ProductRecipeItem{
			BaseModel: domain.BaseModel{
				CreatedAt: now,
				CreatedBy: req.Username,
			},
			IngredientId:   itemReq.IngredientId,
			Quantity:       &itemReq.Quantity,
			UnitId:         &itemReq.UnitId,
			IngredientCost: ingredientCost,
		}

		items = append(items, newItem)
		itemsWithDetail = append(itemsWithDetail, recipeItemWithDetail{
			Item:       newItem,
			Ingredient: *ingredient,
			Unit:       *recipeUnit,
		})
	}

	var costPerUnit float64
	if req.EstimatedYield > 0 {
		costPerUnit = totalCost / req.EstimatedYield
	}

	newRecipe.TotalCost = totalCost
	newRecipe.CostPerUnit = costPerUnit

	if err := p.repoItem.CreateWithItems(newRecipe, items); err != nil {
		return nil, err
	}

	for i := range itemsWithDetail {
		itemsWithDetail[i].Item.Id = items[i].Id
	}

	product.EstimatedCost = &totalCost
	if err := p.repo.Update(product); err != nil {
		return nil, err
	}

	response := p.mappingRecipeResponse(*newRecipe, itemsWithDetail, *yieldUnit)
	return &response, nil
}

func (u *ProductUsecase) validateUom(unitId string) (*domain.Uom, error) {
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

func (u *ProductUsecase) validationCategory(categoyId string) (*domain.Categories, error) {
	category, err := u.repoCategory.FindById(categoyId)
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
func (u *ProductUsecase) validationIngredient(ingredientId string) (*domain.Ingredient, error) {
	ingredient, err := u.repoIngredient.FindById(ingredientId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &middleware.GlobalError{
				Message: "Ingredient not found",
				Code:    middleware.ErrNotFound.Code,
				Status:  middleware.ErrNotFound.Status,
			}
		}

		return nil, err
	}
	return ingredient, nil
}

func (u *ProductUsecase) mappingResponse(product []domain.Product) []domain.ProductResponse {
	result := make([]domain.ProductResponse, 0, len(product))

	for _, product := range product {
		result = append(result, u.mapppingSingleResponse(product))
	}
	return result
}

func (u *ProductUsecase) mapppingSingleResponse(product domain.Product) domain.ProductResponse {
	return domain.ProductResponse{
		Id:           product.Id,
		Code:         product.Code,
		Name:         product.Name,
		CategoryID:   product.Category.Id,
		CategoryName: product.Category.Name,
		UnitID:       product.Unit.Id,
		UnitSymbol:   product.Unit.Symbol,
		Type:         product.Type,
		TargetMargin: product.TargetMargin,
		StockQty:     product.StockQty,
		IsActive:     product.IsActive,
	}
}

func (p *ProductUsecase) mappingRecipeResponse(
	recipe domain.ProductRecipes,
	itemsWithDetail []recipeItemWithDetail,
	yieldUnit domain.Uom,
) domain.RecipeResponse {
	itemResponses := make([]domain.ProductRecipeItemResponse, 0, len(itemsWithDetail))

	for _, detail := range itemsWithDetail {
		quantity := 0.0
		if detail.Item.Quantity != nil {
			quantity = *detail.Item.Quantity
		}

		itemResponses = append(itemResponses, domain.ProductRecipeItemResponse{
			Id:             detail.Item.Id,
			IngredientId:   detail.Item.IngredientId,
			IngredientName: detail.Ingredient.Name,
			UnitId:         detail.Unit.Id,
			UnitSymbol:     detail.Unit.Symbol,
			Quantity:       quantity,
		})
	}

	estimatedYield := 0.0
	if recipe.EstimatedYield != nil {
		estimatedYield = *recipe.EstimatedYield
	}

	return domain.RecipeResponse{
		Id:              recipe.Id,
		VersionNumber:   recipe.VersionNumber,
		IsActive:        recipe.IsActive,
		Notes:           recipe.Notes,
		Items:           itemResponses,
		EstimatedCost:   recipe.TotalCost,
		EstimatedYield:  estimatedYield,
		YieldUnitName:   yieldUnit.Name,
		YieldUnitSymbol: yieldUnit.Symbol,
		CostPerUnit:     recipe.CostPerUnit,
	}
}
