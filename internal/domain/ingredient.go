package domain

type Ingredient struct {
	BaseModelWithDeleted
	TenantId      *string  `gorm:"type:uuid"`
	CategoryId    string   `gorm:"type:uuid"`
	UnitId        string   `gorm:"type:uuid"`
	Name          string   `gorm:"type:varchar(255);not null;index"`
	StockQuantity float64  `gorm:"type:numeric;not null"`
	MinimumStock  float64  `gorm:"type:numeric;not null"`
	PricePerUnit  *float64 `gorm:"type:numeric;default:null"`
	PurchasePrice float64  `gorm:"type:numeric;not null"`
	AveragePrice  float64  `gorm:"type:numeric;not null"`
	IsActive      bool     `gorm:"type:bool;default:true;not null"`
}

type IngredientRequest struct {
	TenantId      *string  `json:"tenant_id,omitempty"`
	CategoryId    string   `json:"category_id" binding:"required"`
	UnitId        string   `json:"unit_id" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	StockQuantity *float64  `json:"stock_qty"`
	MinimumStock  float64  `json:"min_stock" binding:"required"`
	PricePerUnit  *float64 `json:"price_per_unit"`
	PurchasePrice float64  `json:"purchase_price"`
	AveragePrice  float64  `json:"average_price"`
	IsActive      bool     `json:"is_active"`
	Username      string
}

type IngredientRequestUpdate struct {
	Id string `json:"id" binding:"required"`
	IngredientRequest
}

type IngredientResponse struct {
	Id            *string  `json:"id"`
	CategoryId    string   `json:"category_id" `
	UnitId        string   `json:"unit_id"`
	Name          string   `json:"name"`
	StockQuantity float64  `json:"stock_qty"`
	MinimumStock  *float64 `json:"min_stock"`
	PricePerUnit  *float64 `json:"price_per_unit"`
	PurchasePrice float64  `json:"purchase_price"`
	AveragePrice  float64  `json:"average_price"`
	IsActive      bool     `json:"is_active"`
}
