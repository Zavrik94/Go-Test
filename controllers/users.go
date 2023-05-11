package controllers

import (
	"github.com/gin-gonic/gin"
	"go-test/database/models"
	"net/http"
	"strconv"
)

func Profile(c *gin.Context) {
	id, exists := c.Get("userID")

	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	user.FindByID(id.(int))

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
	})
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.RoleID = 1

	if !user.Create() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func CreateAdmin(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.RoleID = 2

	if !user.Create() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func DeleteUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))

	var user models.User

	user.FindByID(userId)

	user.SoftDelete()
}

type UserJSON struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func UsersList(c *gin.Context) {
	users := models.GetAllUsers()

	usersJSON := make([]UserJSON, len(users))

	for i, user := range users {
		usersJSON[i] = UserJSON{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role.Name,
		}
	}

	// Return the JSON response
	c.JSON(http.StatusOK, gin.H{
		"cars": usersJSON,
	})
}
