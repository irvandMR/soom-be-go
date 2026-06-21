package domain

type ProductRecipeItem struct {
	BaseModel
	RecipeId       string   `gorm:"column:recipe_id;type:uuid;not null"`
	IngredientId   string   `gorm:"column:ingredient_id;type:uuid;not null"`
	Quantity       *float64 `gorm:"column:quantity;type:numeric"`
	UnitId         *string  `gorm:"column:unit_id;type:uuid"`
	IngredientCost float64  `gorm:"column:ingredient_cost;type:numeric;not null"`
}
type ProductRecipeItemRequest struct {
	IngredientId string  `json:"ingredient_id" binding:"required"`
	UnitId       string  `json:"unit_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
}

type ProductRecipeItemResponse struct {
	Id             string  `json:"id"`
	IngredientId   string  `json:"ingredient_id" `
	IngredientName string  `json:"ingredient_name" `
	UnitId         string  `json:"unit_id"`
	UnitSymbol     string  `json:"unit_symbol"`
	Quantity       float64 `json:"quantity"`
}
