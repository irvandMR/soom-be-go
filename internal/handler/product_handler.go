package handler

import (
	"net/http"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/repository"
	"soom-be-go/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	repoProduct := repository.NewProductRepository(db)
	repoCategory := repository.NewCategoriesRepository(db)
	repoUom := repository.NewUomRepository(db)
	repoTenant := repository.NewTenantRepository(db)
	repoRecipe := repository.NewProductRecipeRepository(db)
	repoIngredient := repository.NewIngredientRepository(db)
	repoItem := repository.NewProductRecipeItemRepository(db)
	uc := usecase.NewProductUsecase(repoProduct, repoCategory, repoUom, repoTenant, repoRecipe, repoIngredient, repoItem)
	return &ProductHandler{
		usecase: uc,
	}
}

func (p *ProductHandler) GetAll(c *gin.Context) {
	var req domain.ProductQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", err.Error()))
		return
	}

	tenantId := c.GetString("trenatId")
	if tenantId != "" {
		req.TenantId = &tenantId
	}

	data, err := p.usecase.GetAllProducts(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse("Success get products", data))
}

func (p *ProductHandler) CreatedProduct(c *gin.Context) {
	var req domain.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	req.Username = c.GetString("username")
	tenantId := c.GetString("tenantId")
	if tenantId != "" {
		req.TenantId = &tenantId
	}
	data, err := p.usecase.CreateProduct(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success create products", data))

}
func (p *ProductHandler) UpdateProduct(c *gin.Context) {
	var req domain.ProductRequestUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	req.Username = c.GetString("username")
	tenantId := c.GetString("tenantId")
	if tenantId != "" {
		req.TenantId = &tenantId
	}
	data, err := p.usecase.UpdateProduct(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success update products", data))

}

func (p *ProductHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	data, err := p.usecase.GetProductById(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success get product", data))
}

func (p *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	user := c.GetString("username")
	if err := p.usecase.DeleteProduct(id, user); err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, SuccessResponse("Success Delete product", nil))
}

func (p *ProductHandler) CreatedProductRecipe(c *gin.Context) {
	var req domain.ProductRecipesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request", handleValidationError(err)))
		return
	}
	req.Username = c.GetString("username")
	req.TenantId = c.GetString("tenantId")
	data, err := p.usecase.SaveProductRecipe(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse("Success create products", data))

}
