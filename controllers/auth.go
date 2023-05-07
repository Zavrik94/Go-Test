package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-test/database"
	"go-test/database/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	db := database.GetDB()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login details123"})
		return
	}

	password := user.Password

	// Check if user exists in database
	if err := db.Unscoped().Where("deleted_at IS NULL").Where("email = ?", user.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login details"})
		return
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	tokenDb := models.Token{
		UserID:    &user.ID,
		ExpiredAt: time.Now().Local().Add(time.Hour * time.Duration(24)),
	}
	if !tokenDb.Save() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not generate token"})
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Logout(c *gin.Context) {
	userID, exist := c.Get("userID")

	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}

	models.DeleteToken(userID.(uint))

	c.JSON(http.StatusNoContent, gin.H{})
}
