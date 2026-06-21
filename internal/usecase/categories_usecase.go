package usecase

import (
	"errors"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/middleware"
	"soom-be-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type CategoriesUsecase struct {
	repo repository.CategoriesRepository
}

func NewCategoriesUsecase(r repository.CategoriesRepository) *CategoriesUsecase {
	return &CategoriesUsecase{repo: r}
}

func (u *CategoriesUsecase) GetAllCategories(req domain.CategoriesQueryRequest) (*domain.PaginationResponse, error) {
	req.Normalize()

	categories, total, err := u.repo.FindAll(req)
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
		Data:       categories,
	}, nil
}
func (u *CategoriesUsecase) GetAllCategoriesWithoutPagination(tenantId *string) ([]domain.CategoriesResponse, error) {
	categories, err := u.repo.FindTypeByTenant(tenantId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}

	var mappingResponse []domain.CategoriesResponse
	for _, categorie := range categories {
		mappingResponse = append(mappingResponse, domain.CategoriesResponse{
			Id:       categorie.Id,
			Code:     categorie.Code,
			Name:     categorie.Name,
			Type:     categorie.Type,
			IsActive: categorie.IsActive,
		})
	}
	return mappingResponse, err
}

func (u *CategoriesUsecase) GetCategoriesById(id string) (*domain.Categories, error) {
	categories, err := u.repo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	return categories, nil
}

func (u *CategoriesUsecase) GetCategoriesByType(tenantId *string) ([]domain.Categories, error) {
	categories, err := u.repo.FindTypeByTenant(tenantId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	return categories, nil
}

func (u *CategoriesUsecase) CreateCategories(req domain.CategoriesRequest) (*domain.Categories, error) {
	categories := &domain.Categories{
		BaseModelWithDeleted: domain.BaseModelWithDeleted{
			BaseModel: domain.BaseModel{
				CreatedAt: time.Now(),
				CreatedBy: req.Username,
			},
		},
		TenantId: req.TenantId,
		Code:     req.Code,
		Name:     req.Name,
		IsActive: req.IsActive,
		Type:     req.Type,
	}
	err := u.repo.Create(categories)
	return categories, err
}

func (u *CategoriesUsecase) UpdateCategories(req domain.CategoriesRequestUpdate) (*domain.Categories, error) {
	categories, err := u.repo.FindById(req.Id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	categories.Code = req.Code
	categories.Name = req.Name
	categories.IsActive = req.IsActive
	categories.UpdatedAt = &now
	categories.UpdatedBy = &req.Username

	err = u.repo.Update(categories)
	return categories, err
}

func (u *CategoriesUsecase) DeleteCategories(id string, deletedBy string) error {
	_, err := u.repo.FindById(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(id, deletedBy)
}
