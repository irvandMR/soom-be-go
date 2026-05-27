package main

import (
	"os"
	"soom-be-go/internal/config"
	"soom-be-go/internal/handler"
	"soom-be-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db, sqlDB := config.InitDB()
	_ = db

	config.RunMigrations(sqlDB)

	port := os.Getenv("APP_PORT")

	r := gin.Default()

	r.Use(middleware.Errorhandler())

	handler.RegisterRoutes(r, db)

	r.Run(":" + port)
}
