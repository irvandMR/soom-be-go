package domain

type Product struct {
	BaseModelWithDeleted
	Code          string   `gorm:"type:varchar(100);unique;not null"`
	CategoryID    string   `gorm:"type:uuid;not null"`
	UnitID        string   `gorm:"type:uuid;not null"`
	Name          string   `gorm:"type:varchar(255);not null"`
	Type          string   `gorm:"type:varchar(100);not null"`
	DefaultPrice  *float64 `gorm:"type:numeric"`
	StockQty      *float64 `gorm:"type:numeric"`
	EstimatedCost *float64 `gorm:"type:numeric"`
	TargetMargin  float64  `gorm:"type:numeric;not null"`
	TenantID      *string  `gorm:"type:uuid;not null"`
	IsActive      bool     `gorm:"default:false;not null"`

	Category Categories `gorm:"foreignKey:CategoryId"`
	Unit     Uom        `gorm:"foreignKey:UnitId"`
}

type ProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	Code       string  `json:"code" binding:"required"`
	CategoryID string  `json:"category_id" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	UnitID     string  `json:"unit_id" binding:"required"`
	IsActive   bool    `json:"is_active"`
	TenantId   *string `json:"tenant_id,omitempty"`
	Username   string
}

type ProductRequestUpdate struct {
	Id string `json:"id" binding:"required"`
	ProductRequest
}

type ProductQueryRequest struct {
	PaginationRequest
	Search     string  `form:"search"`
	CategoryId string  `form:"category_id"`
	Type       string  `form:"type"`
	TenantId   *string `form:"-"`
}

type ProductResponse struct {
	Id            string   `json:"id"`
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	CategoryID    string   `json:"category_id"`
	CategoryName  string   `json:"category_name"`
	UnitID        string   `json:"unit_id"`
	UnitSymbol    string   `json:"unit_symbol"`
	DefaultPrice  *float64 `json:"default_price"`
	StockQty      *float64 `json:"stock_qty"`
	EstimatedCost *float64 `json:"estimated_cost"`
	TargetMargin  float64  `json:"target_margin"`
	IsActive      bool     `json:"is_active"`
}
