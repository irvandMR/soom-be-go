package domain

type Categories struct {
	BaseModelWithDeleted
	TenantId *string `gorm:"type:uuid"`
	Code     string  `gorm:"size:100;not null;unique"`
	Name     string  `gorm:"size:100;not null"`
	Type     string  `gorm:"size:50;not null"`
	IsActive bool    `gorm:"default:false;not null"`
}

type CategoriesRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	IsActive bool   `json:"is_active"`
	Username string
	TenantId *string
}

type CategoriesRequestUpdate struct {
	Id string `json:"id" binding:"required"`
	CategoriesRequest
}

type CategoriesResponse struct {
	Id       string `json:"id"`
	Code     string `json:"code" `
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	Type     string `json:"type"`
}

type CategoriesQueryRequest struct {
	PaginationRequest
	Search string `form:"search"`
	Type   string `from:"type"`
}
