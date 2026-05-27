package handler

import (
	"net/http"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/repository"
	"soom-be-go/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UomHandler struct {
	usecase *usecase.UomUsecase
}

func NewUomHandler(db *gorm.DB) *UomHandler {
	repo := repository.NewUomRepository(db)
	uc := usecase.NewUomUsecase(repo)
	return &UomHandler{
		usecase: uc,
	}
}

func (h *UomHandler) GetAll(c *gin.Context) {
	var req domain.UomQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", err.Error()))
		return
	}
	data, err := h.usecase.GetAll(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success get uom", data))
}

func (h *UomHandler) GetUomById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.usecase.GetUomById(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success get uom", data))
}

func (h *UomHandler) CreateUom(c *gin.Context) {
	var req domain.UomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	data, err := h.usecase.CreateUom(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success create uom", data))
}

func (h *UomHandler) UpdateUom(c *gin.Context) {
	var req domain.UomRequestUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	user := c.GetString("username")
	data, err := h.usecase.UpdateUom(req, user) // Replace "system" with the actual user ID if available
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success update uom", data))
}

func (h *UomHandler) DeleteUom(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteUom(id, "system"); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse("Success delete uom", nil))
}
