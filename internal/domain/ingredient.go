package domain

type Ingredient struct {
	BaseModelWithDeleted
	TenantId      *string  `gorm:"type:uuid"`
	CategoryId    string   `gorm:"type:uuid"`
	UnitId        string   `gorm:"type:uuid"`
	Name          string   `gorm:"type:varchar(255);not null;index"`
	StockQuantity *float64 `gorm:"column:stock_qty;type:numeric;"`
	MinimumStock  float64  `gorm:"column:min_stock_qty;type:numeric;not null"`
	PricePerUnit  *float64 `gorm:"type:numeric;default:null"`
	PurchasePrice *float64 `gorm:"type:numeric;"`
	AveragePrice  *float64 `gorm:"type:numeric;"`
	IsActive      bool     `gorm:"type:bool;default:true;not null"`
	Status        string   `gorm:"type:varchar(50);index"`

	Category Categories `gorm:"foreignKey:CategoryId"`
	Unit     Uom        `gorm:"foreignKey:UnitId"`
}

type IngredientRequest struct {
	TenantId     *string `json:"tenant_id,omitempty"`
	CategoryId   string  `json:"category_id" binding:"required"`
	UnitId       string  `json:"unit_id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	MinimumStock float64 `json:"min_stock" binding:"required"`
	IsActive     bool    `json:"is_active"`
	Username     string
}

type IngredientRequestUpdate struct {
	Id            string  `json:"id" binding:"required"`
	PurchasePrice float64 `json:"purchase_price"`
	HistoryId     string  `json:"history_id,omitempty"`
	IngredientRequest
}

type IngredientResponse struct {
	Id            string   `json:"id"`
	CategoryId    string   `json:"category_id"`
	CategoryName  string   `json:"category_name"`
	UnitId        string   `json:"unit_id"`
	UnitSymbol    string   `json:"unit_symbol"`
	Name          string   `json:"name"`
	StockQuantity *float64 `json:"stock_qty"`
	MinimumStock  float64  `json:"min_stock"`
	PricePerUnit  *float64 `json:"price_per_unit"`
	PurchasePrice *float64 `json:"purchase_price"`
	AveragePrice  *float64 `json:"average_price"`
	IsActive      bool     `json:"is_active"`
	Status        string   `json:"status"`
}

type IngredientQueryRequest struct {
	PaginationRequest
	Search       string  `form:"search"`
	CategoriesId string  `form:"category_id"`
	Status       string  `form:"status"`
	TenantId     *string `form:"-"`
}
