package main

import (
	"os"
	"soom-be-go/internal/config"
	"soom-be-go/internal/handler"
	"soom-be-go/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db, sqlDB := config.InitDB()
	_ = db

	config.RunMigrations(sqlDB)
	config.InitRedis()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Tenant-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Errorhandler())

	handler.RegisterRoutes(r, db)

	r.Run(":" + port)
}
