package handler

import (
	"soom-be-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/health", HealthHandler)

	v1 := router.Group("/api/v1")
	{
		// AUTH
		authHandler := NewAuthHandler(db)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}
		// UOM routes
		uomHandler := NewUomHandler(db)
		uom := v1.Group("/uoms")
		uom.Use(middleware.JwtAuth())
		{
			uom.GET("", uomHandler.GetAll)
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
			ctg.GET("/:id", categoriesHandler.GetCategoriesById)
			ctg.POST("", categoriesHandler.CreateCategories)
			ctg.POST("/update", categoriesHandler.UpdateCategories)
			ctg.DELETE(":id", categoriesHandler.DeleteCategories)
		}

		// Ingredient
		ingredientHandle := NewIngredientHandler(db)
		ing := v1.Group("/ingredient")
		ing.Use(middleware.JwtAuth())
		{
			ing.GET("", ingredientHandle.GetAll)
			ing.GET("/:id", ingredientHandle.GetIngredientById)
			ing.POST("", ingredientHandle.CreateIngredient)
			ing.POST("/update", ingredientHandle.UpdateIngredient)
			ing.DELETE(":id", ingredientHandle.DeleteIngredient)
		}
	}

}
