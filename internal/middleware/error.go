package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GlobalError struct {
	Message string
	Code    string
	Status  int
}

func (e *GlobalError) Error() string {
	return e.Message
}

// Error constants
var (
	ErrNotFound       = &GlobalError{Code: "NOT_FOUND", Message: "Data not found", Status: http.StatusNotFound}
	ErrBadRequest     = &GlobalError{Code: "BAD_REQUEST", Message: "Invalid request", Status: http.StatusBadRequest}
	ErrInternalServer = &GlobalError{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error", Status: http.StatusInternalServerError}
)

func Errorhandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		// Custom Global Error
		var globalErr *GlobalError
		if errors.As(err, &globalErr) {
			c.JSON(globalErr.Status, gin.H{
				"status":  "error",
				"code":    globalErr.Code,
				"message": globalErr.Message,
			})
			return
		}

		// GORM error — record not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"code":    "NOT_FOUND",
				"message": "Data not found",
			})
			return
		}

		// Default
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Internal server error",
		})
	}
}
