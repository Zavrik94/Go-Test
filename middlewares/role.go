package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleMiddleware(requiredRole ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		for _, role := range requiredRole {
			if role == userRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are not authorized to access this resource",
		})
	}
}
