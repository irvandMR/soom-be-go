package domain

import "time"

type IngredientsStockHistory struct {
	BaseModel
	IngredientId  string  `gorm:"type:uuid;not null;index"`
	Type          string  `gorm:"type:varchar(25);not null"` // 'IN', 'OUT', 'ADJUSTMENT', 'WASTE'
	Quantity      float64 `gorm:"type:numeric;not null"`
	PurchasePrice float64 `gorm:"type:numeric;not null;default:0"`
	Notes         *string `gorm:"type:varchar(255);default:null"`
	ReferenceId   *string `gorm:"type:varchar(255);default:null"`
	ReferenceType *string `gorm:"type:varchar(100);default:null"`
}

type StockInRequest struct {
	IngredientId  string  `json:"ingredient_id" binding:"required"`
	Quantity      float64 `json:"quantity" binding:"required,gt=0"`
	PurchasePrice float64 `json:"purchase_price" binding:"required,gt=0"`
	Notes         *string `json:"notes"`
	Username      string
}

type StockOutRequest struct {
	IngredientId string  `json:"ingredient_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	Notes        *string `json:"notes"`
}

type IngredientsStockHistoryRequest struct {
	IngredientId string     `form:"ingredient_id" uri:"id"`
	StartDate    *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate      *time.Time `form:"end_date" time_format:"2006-01-02"`
}

type IngredientsStockHistoryResponse struct {
	Id            string    `json:"id"`
	Type          string    `json:"type"`
	Quantity      float64   `json:"quantity"`
	PurchasePrice float64   `json:"purchase_price"`
	Notes         *string   `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}
