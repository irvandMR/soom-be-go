package domain

type ProductRecipes struct {
	BaseModel
	ProductId      string   `gorm:"column:product_id;type:uuid;not null"`
	VersionNumber  int32    `gorm:"column:version_number;type:int4;not null"`
	Notes          *string  `gorm:"column:notes;type:text"`
	EstimatedYield *float64 `gorm:"column:estimated_yield;type:numeric"`
	UnitId         *string  `gorm:"column:unit_id;type:uuid"`
	IsActive       bool     `gorm:"column:is_active;not null;default:false"`
	TotalCost      float64  `gorm:"column:total_cost;type:numeric;not null;default:0"`
	CostPerUnit    float64  `gorm:"column:cost_per_unit;type:numeric;not null;default:0"`

	Product Product `gorm:"foreignKey:ProductId"`
}

type ProductRecipesRequest struct {
	ProductId      string                     `json:"product_id" binding:"required"`
	EstimatedYield float64                    `json:"estimated_yield" binding:"required,gt=0"`
	Notes          *string                    `json:"notes"`
	Items          []ProductRecipeItemRequest `json:"items" binding:"required,min=1,dive"`
	Username       string                     `json:"-"`
	TenantId       string                     `json:"-"`
	UnitId         string                     `json:"unit_id" binding:"required,uuid"`
}

type RecipeResponse struct {
	Id              string                      `json:"id"`
	VersionNumber   int32                       `json:"version_number"`
	IsActive        bool                        `json:"is_active"`
	Notes           *string                     `json:"notes"`
	Items           []ProductRecipeItemResponse `json:"items"`
	EstimatedCost   float64                     `json:"estimated_cost"`
	EstimatedYield  float64                     `json:"estimated_yield"`
	YieldUnitName   string                      `json:"yield_unit_name"`
	YieldUnitSymbol string                      `json:"yield_unit_symbol"`
	CostPerUnit     float64                     `json:"cost_per_unit"`
}
