package handler

import (
	"soom-be-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/health", HealthHandler(db))

	v1 := router.Group("/api/v1")
	{
		// AUTH
		authHandler := NewAuthHandler(db)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", middleware.JwtAuth(), authHandler.Logout)
			auth.GET("/me", middleware.JwtAuth(), authHandler.me)
		}
		// UOM routes
		uomHandler := NewUomHandler(db)
		uom := v1.Group("/uoms")
		uom.Use(middleware.JwtAuth())
		{
			uom.GET("", uomHandler.GetAll)
			uom.GET("/all", uomHandler.GetUomAll)
			uom.GET("/detail/:id", uomHandler.GetUomById)
			uom.POST("", uomHandler.CreateUom)
			uom.POST("/update", uomHandler.UpdateUom)
			uom.DELETE("/:id", uomHandler.DeleteUom)
		}

		// Categories
		categoriesHandler := NewCategoriesHandler(db)
		ctg := v1.Group("/categories")
		ctg.Use(middleware.JwtAuth())
		{
			ctg.GET("", categoriesHandler.GetAll)
			ctg.GET("/all", categoriesHandler.GetAllCategories)
			ctg.GET("/:id", categoriesHandler.GetCategoriesById)
			ctg.POST("", categoriesHandler.CreateCategories)
			ctg.POST("/update", categoriesHandler.UpdateCategories)
			ctg.DELETE("/:id", categoriesHandler.DeleteCategories)
			ctg.GET("/types", categoriesHandler.GetCategoriesType)
		}

		// Ingredient
		ingredientHandle := NewIngredientHandler(db)
		ing := v1.Group("/ingredient")
		ing.Use(middleware.JwtAuth())
		{
			ing.GET("", ingredientHandle.GetAll)
			ing.GET("/all", ingredientHandle.GetAllIngredient)
			ing.GET("/:id", ingredientHandle.GetIngredientById)
			ing.POST("", ingredientHandle.CreateIngredient)
			ing.POST("/update", ingredientHandle.UpdateIngredient)
			ing.DELETE("/:id", ingredientHandle.DeleteIngredient)

			ing.POST("/stock-in", ingredientHandle.StockIn)
			ing.GET("/history/:id", ingredientHandle.GetHistory)
		}

		// Product
		productHandle := NewProductHandler(db)
		product := v1.Group("/product")
		product.Use(middleware.JwtAuth())
		{
			product.GET("", productHandle.GetAll)
			product.POST("", productHandle.CreatedProduct)
			product.GET("/:id", productHandle.GetProductById)
			product.POST("/update", productHandle.UpdateProduct)
			product.DELETE("/:id", productHandle.DeleteProduct)

			// Recipe

		}
	}

}
