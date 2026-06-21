package handler

import (
	"context"
	"net/http"
	"soom-be-go/internal/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "ok"
		dbStatus := "up"
		redisStatus := "up"

		// Check DB connection
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
			status = "error"
		}

		// Check Redis connection
		if config.RedisClient == nil || config.RedisClient.Ping(context.Background()).Err() != nil {
			redisStatus = "down"
			status = "error"
		}

		httpStatus := http.StatusOK
		if status != "ok" {
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, gin.H{
			"status": status,
			"dependencies": gin.H{
				"database": dbStatus,
				"redis":    redisStatus,
			},
		})
	}
}