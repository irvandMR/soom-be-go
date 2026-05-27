package domain

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
