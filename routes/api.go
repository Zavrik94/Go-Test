package routes

import "github.com/gin-gonic/gin"

func RegisterAPIRoutes(router *gin.Engine) {
	router.GET("/api/users", getUsers)
}

func getUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "getUsers handler123",
	})
}
