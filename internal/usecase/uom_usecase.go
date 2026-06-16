package usecase

import (
	"errors"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/middleware"
	"soom-be-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type UomUsecase struct {
	repo repository.UomRepository
}

func NewUomUsecase(r repository.UomRepository) *UomUsecase {
	return &UomUsecase{repo: r}
}

func (u *UomUsecase) GetAll(req domain.UomQueryRequest) (*domain.PaginationResponse, error) {
	req.Normalize()

	uom, total, err := u.repo.FindAll(req)
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
		Data:       u.mappingResponse(uom),
	}, nil
}

func (u *UomUsecase) GetUomById(id string) (*domain.UomResponse, error) {
	uom, err := u.repo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	result := u.mappingResponse([]domain.Uom{*uom})
	return &result[0], nil
}
func (u *UomUsecase) GetUomAll() ([]domain.UomResponse, error) {
	uoms, err := u.repo.FindAllNoPagination()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound
		}
		return nil, err
	}
	result := u.mappingResponse(uoms)
	return result, nil
}

func (u *UomUsecase) CreateUom(req domain.UomRequest) (*domain.Uom, error) {
	uom := &domain.Uom{
		BaseModelWithDeleted: domain.BaseModelWithDeleted{
			BaseModel: domain.BaseModel{
				CreatedAt: time.Now(),
				CreatedBy: "system",
			},
		},
		Code:             req.Code,
		Name:             req.Name,
		Symbol:           req.Symbol,
		HaveConversion:   req.HaveConversion,
		BaseUnit:         req.BaseUnit,
		ConversionFactor: req.ConversionFactor,
	}
	err := u.repo.Create(uom)
	return uom, err
}

func (u *UomUsecase) UpdateUom(req domain.UomRequestUpdate, updatedBy string) (*domain.Uom, error) {
	uom, err := u.repo.FindById(req.Id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	uom.Code = req.Code
	uom.Name = req.Name
	uom.Symbol = req.Symbol
	uom.HaveConversion = req.HaveConversion
	uom.BaseUnit = req.BaseUnit
	uom.ConversionFactor = req.ConversionFactor
	uom.UpdatedAt = &now
	uom.UpdatedBy = &updatedBy
	err = u.repo.Update(uom)
	return uom, err
}

func (u *UomUsecase) DeleteUom(id string, deletedBy string) error {
	_, err := u.repo.FindById(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(id, deletedBy)
}

func (u *UomUsecase) mappingResponse(uoms []domain.Uom) []domain.UomResponse {
	result := make([]domain.UomResponse, 0, len(uoms))

	for _, uom := range uoms {
		result = append(result, domain.UomResponse{
			Id:               uom.Id,
			Code:             uom.Code,
			Name:             uom.Name,
			HaveConversion:   uom.HaveConversion,
			BaseUnit:         uom.BaseUnit,
			Symbol:           uom.Symbol,
			ConversionFactor: uom.ConversionFactor,
		})
	}
	return result
}
