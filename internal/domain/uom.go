package domain

type Uom struct {
	BaseModelWithDeleted
	Code             string  `gorm:"size:50;unique;not null"`
	Name             string  `gorm:"size:100;not null"`
	Symbol           string  `gorm:"size:100;not null"`
	HaveConversion   bool    `gorm:"default:false;not null"`
	BaseUnit         *string `gorm:"size:25"`
	ConversionFactor *float64
}

type UomRequest struct {
	Code             string   `json:"code" binding:"required"`
	Name             string   `json:"name" binding:"required"`
	Symbol           string   `json:"symbol" binding:"required"`
	HaveConversion   bool     `json:"have_conversion"`
	BaseUnit         *string  `json:"base_unit,omitempty"`
	ConversionFactor *float64 `json:"conversion_factor,omitempty"`
}

type UomRequestUpdate struct {
	Id string `json:"id" binding:"required"`
	UomRequest
}

type UomResponse struct {
	Id               string   `json:"id"`
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	Symbol           string   `json:"symbol"`
	HaveConversion   bool     `json:"have_conversion"`
	BaseUnit         *string  `json:"base_unit,omitempty"`
	ConversionFactor *float64 `json:"conversion_factor,omitempty"`
}

type UomQueryRequest struct {
	PaginationRequest
	Search   string `form:"search"`
	Baseunit string `form:"base_unit"`
}
