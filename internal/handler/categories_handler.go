package handler

import (
	"net/http"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/repository"
	"soom-be-go/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoriesHandler struct {
	usecase *usecase.CategoriesUsecase
}

func NewCategoriesHandler(db *gorm.DB) *CategoriesHandler {
	repo := repository.NewCategoriesRepository(db)
	uc := usecase.NewCategoriesUsecase(repo)
	return &CategoriesHandler{usecase: uc}
}

func (h *CategoriesHandler) GetAll(c *gin.Context) {
	var req domain.CategoriesQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", err.Error()))
		return
	}
	tenantId := c.GetString("tenantId")
	if tenantId != "" {
		req.TenantId = &tenantId
	}
	data, err := h.usecase.GetAllCategories(req)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, SuccessResponse("Success get categories", data))
}

func (h *CategoriesHandler) GetCategoriesById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.usecase.GetCategoriesById(id)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, SuccessResponse("Success get categories", data))
}

func (h *CategoriesHandler) GetCategoriesType(c *gin.Context) {
	tenantIdStr := c.GetString("tenantId")

	var tenantId *string
	if tenantIdStr != "" {
		tenantId = &tenantIdStr
	}

	data, err := h.usecase.GetCategoriesByType(tenantId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse("Success get categorie types", data))
}
func (h *CategoriesHandler) GetAllCategories(c *gin.Context) {
	tenantIdStr := c.GetString("tenantId")

	var tenantId *string
	if tenantIdStr != "" {
		tenantId = &tenantIdStr
	}

	data, err := h.usecase.GetAllCategoriesWithoutPagination(tenantId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse("Success get categorie types", data))
}

func (h *CategoriesHandler) CreateCategories(c *gin.Context) {
	var req domain.CategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	req.Username = c.GetString("username")
	tenantId := c.GetString("tenantId")
	if tenantId != "" {
		req.TenantId = &tenantId
	}
	data, err := h.usecase.CreateCategories(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success create uom", data))
}

func (h *CategoriesHandler) UpdateCategories(c *gin.Context) {
	var req domain.CategoriesRequestUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	req.Username = c.GetString("username")
	tenantId := c.GetString("tenantId")
	if tenantId != "" {
		req.TenantId = &tenantId
	}
	data, err := h.usecase.UpdateCategories(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success update uom", data))
}

func (h *CategoriesHandler) DeleteCategories(c *gin.Context) {
	id := c.Param("id")
	user := c.GetString("username")
	if err := h.usecase.DeleteCategories(id, user); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse("Success delete categories", nil))
}
