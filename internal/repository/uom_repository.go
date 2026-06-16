package repository

import (
	"fmt"
	"soom-be-go/internal/domain"
	"time"

	"gorm.io/gorm"
)

type UomRepository interface {
	Create(uom *domain.Uom) error
	Update(uom *domain.Uom) error
	FindById(id string) (*domain.Uom, error)
	Delete(id string, deletedBy string) error
	FindAll(req domain.UomQueryRequest) ([]domain.Uom, int64, error)
	FindAllNoPagination() ([]domain.Uom, error)
}

type uomRepository struct {
	db *gorm.DB
}

// FindAllNoPagination implements [UomRepository].
func (r *uomRepository) FindAllNoPagination() ([]domain.Uom, error) {
	var uom []domain.Uom
	result := r.db.Where(" deleted_at is null").Find(&uom)
	if result.Error != nil {
		return nil, result.Error
	}
	return uom, nil
}

// FindAll implements [UomRepository].
func (r *uomRepository) FindAll(req domain.UomQueryRequest) ([]domain.Uom, int64, error) {
	var uom []domain.Uom
	var total int64

	fmt.Println("req : ", req)

	query := r.db.Model(&domain.Uom{}).Where("deleted_at is null")

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if req.Baseunit != "" {
		query = query.Where("base_unit = ?", req.Baseunit)
	}

	query.Count(&total)

	result := query.Limit(req.Limit).Offset(req.Offset()).Find(&uom)
	return uom, total, result.Error

}

func NewUomRepository(db *gorm.DB) UomRepository {
	return &uomRepository{
		db: db,
	}
}

func (r *uomRepository) Create(uom *domain.Uom) error {
	result := r.db.Create(uom)
	return result.Error
}

func (r *uomRepository) Update(uom *domain.Uom) error {
	result := r.db.Save(uom)
	return result.Error
}

func (r *uomRepository) FindById(id string) (*domain.Uom, error) {
	var uom domain.Uom
	result := r.db.Where("id = ? AND deleted_at is null", id).First(&uom)
	if result.Error != nil {
		return nil, result.Error
	}
	return &uom, nil
}

func (r *uomRepository) Delete(id string, deletedBy string) error {
	now := time.Now()
	result := r.db.Model(&domain.Uom{}).Where("id = ? AND deleted_at is null", id).Updates(map[string]interface{}{
		"deleted_by": deletedBy,
		"deleted_at": now,
	})
	return result.Error
}
