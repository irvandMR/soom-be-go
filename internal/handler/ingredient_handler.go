package handler

import (
	"net/http"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/repository"
	"soom-be-go/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IngredientHandler struct {
	usecase *usecase.IngredientUsecase
}

func NewIngredientHandler(db *gorm.DB) *IngredientHandler {
	repo := repository.NewIngredientRepository(db)
	uc := usecase.NewIngredientUsecase(repo)
	return &IngredientHandler{usecase: uc}
}

func (h *IngredientHandler) GetAll(c *gin.Context) {
	var req domain.PaginationRequest
	data, err := h.usecase.GetAllIngredient(req)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, SuccessResponse("Success get ingredients", data))
}

func (h *IngredientHandler) GetIngredientById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.usecase.GetIngredientById(id)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, SuccessResponse("Success get ingredients", data))
}

func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
	var req domain.IngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	req.Username = c.GetString("username")
	tenantId := c.GetString("tenantId")
	if tenantId != "" {
		req.TenantId = &tenantId
	}
	data, err := h.usecase.CreateIngredient(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success create ingredients", data))
}

func (h *IngredientHandler) UpdateIngredient(c *gin.Context) {
	var req domain.IngredientRequestUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	req.Username = c.GetString("username")
	tenantId := c.GetString("tenantId")
	if tenantId != "" {
		req.TenantId = &tenantId
	}
	data, err := h.usecase.UpdateIngredient(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success update ingredients", data))
}

func (h *IngredientHandler) DeleteIngredient(c *gin.Context) {
	id := c.Param("id")
	user := c.GetString("username")
	if err := h.usecase.DeleteIngredient(id, user); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse("Success delete categories", nil))
}
