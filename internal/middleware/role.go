package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"code":    "FORBIDDEN",
			"message": "You don't have permission to access this resource",
		})
		c.Abort()
	}
}
