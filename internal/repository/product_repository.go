package repository

import (
	"fmt"
	"soom-be-go/internal/domain"
	"time"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(req *domain.Product) error
	FindAllPagination(req domain.ProductQueryRequest) ([]domain.Product, int64, error)
	FindById(id string) (*domain.Product, error)
	Update(req *domain.Product) error
	Delete(id string, deletedBy string) error
}

type productRepository struct {
	db *gorm.DB
}

// Delete implements [ProductRepository].
func (p *productRepository) Delete(id string, deletedBy string) error {
	now := time.Now()
	result := p.db.Model(&domain.Ingredient{}).Where("id = ? and deleted_at is null", id).Updates(map[string]interface{}{
		"deleted_by": deletedBy,
		"deleted_at": now,
	})
	return result.Error
}

// FindById implements [ProductRepository].
func (p *productRepository) FindById(id string) (*domain.Product, error) {
	var product domain.Product
	result := p.db.Preload("Category").Preload("Unit").Where("id = ? and deleted_at is null", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

// Update implements [ProductRepository].
func (p *productRepository) Update(req *domain.Product) error {
	return p.db.Save(req).Error
}

// Create implements [ProductRepository].
func (p *productRepository) Create(req *domain.Product) error {
	return p.db.Create(req).Error
}

// FindAllPagination implements [ProductRepository].
func (p *productRepository) FindAllPagination(req domain.ProductQueryRequest) ([]domain.Product, int64, error) {
	var product []domain.Product
	var total int64

	fmt.Printf("req: %+v\n", req)

	query := p.db.Model(&domain.Product{}).Where("deleted_at is null")

	if req.TenantId != nil {
		query = query.Where("tenant_id = ?", *req.TenantId)
	}

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if req.CategoryId != "" {
		query = query.Where("category_id = ?", req.CategoryId)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	query.Count(&total)
	result := query.Preload("Category").Preload("Unit").Limit(req.Limit).Offset(req.Offset()).Find(&product)
	return product, total, result.Error

}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
