package controllers

import (
	"github.com/gin-gonic/gin"
	"go-test/database"
	"go-test/database/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Profile(c *gin.Context) {
	db := database.GetDB()

	id, exists := c.Get("userID")

	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	db.Where("id = ?", id).First(&user)

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

	if !createUser(user) {
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

	if !createUser(user) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func createUser(user models.User) bool {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}

	user.Password = string(hashedPassword)
	db := database.GetDB()

	var count int64
	db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return false
	}

	if err := db.Create(&user).Error; err != nil {
		return false
	}

	return true
}

func DeleteUser(c *gin.Context) {
	userId := c.Param("id")

	var user models.User

	db := database.GetDB()
	db.Model(&models.User{}).Where("id = ?", userId).First(&user)
	now := time.Now()
	user.DeletedAt = &now
	db.Save(&user)
}

type UsersJSON struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func UsersList(c *gin.Context) {
	var users []models.User

	db := database.GetDB()
	db.Preload("Role").Model(&models.User{}).Find(&users)

	var usersJSON []UsersJSON
	for _, user := range users {
		carJSON := UsersJSON{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role.Name,
		}
		usersJSON = append(usersJSON, carJSON)
	}

	// Return the JSON response
	c.JSON(http.StatusOK, gin.H{
		"cars": usersJSON,
	})
}
