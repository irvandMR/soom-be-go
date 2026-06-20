package usecase

import (
	"errors"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/middleware"
	"soom-be-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type ProductUsecase struct {
	repo         repository.ProductRepository
	repoCategory repository.CategoriesRepository
	repoUom      repository.UomRepository
}

func NewProductUsecase(repoProduct repository.ProductRepository, repoCategory repository.CategoriesRepository, repoUom repository.UomRepository) *ProductUsecase {
	return &ProductUsecase{repo: repoProduct, repoCategory: repoCategory, repoUom: repoUom}
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

	product := &domain.Product{
		BaseModelWithDeleted: domain.BaseModelWithDeleted{
			BaseModel: domain.BaseModel{
				CreatedAt: time.Now(),
				CreatedBy: req.Username,
			},
		},
		CategoryID: category.Id,
		UnitID:     uom.Id,
		Code:       req.Code,
		Name:       req.Name,
		Type:       req.Type,
		TenantID:   req.TenantId,
		IsActive:   req.IsActive,
	}

	err = p.repo.Create(product)

	response := p.mapppingSingleResponse(*product)
	return &response, err
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
		CategoryName: product.Category.Name,
		UnitSymbol:   product.Unit.Symbol,
		Type:         product.Type,
		TargetMargin: product.TargetMargin,
		StockQty:     product.StockQty,
		IsActive:     product.IsActive,
	}
}
