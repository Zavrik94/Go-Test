package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-test/database"
	"go-test/database/models"
	"go-test/enums"
	"go-test/middlewares"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UsersJSON struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func RegisterAuthRoutes(router *gin.Engine) {
	router.POST("/auth/login", login)
	router.POST("/auth/logout", logout)

	router.POST("/register", register)

	router.GET("/profile", middlewares.AuthMiddleware(), profile)
	router.GET("/users", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.Admin()), usersList)
	router.POST("/admin", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.SuperAdmin()), createAdmin)
	router.POST("/delete/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.Admin()), deleteUser)
}

func login(c *gin.Context) {
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

func logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "logout handler",
	})
}

func profile(c *gin.Context) {
	db := database.GetDB()

	id, exists := c.Get("userID")

	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	db.Where("id = ?", id).First(&user)

	c.JSON(200, gin.H{
		"email": user.Email,
	})
}

func register(c *gin.Context) {
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

	c.JSON(204, gin.H{})
}

func createAdmin(c *gin.Context) {
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

	c.JSON(204, gin.H{})
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

func deleteUser(c *gin.Context) {
	userId := c.Param("id")

	var user models.User

	db := database.GetDB()
	db.Model(&models.User{}).Where("id = ?", userId).First(&user)
	now := time.Now()
	user.DeletedAt = &now
	db.Save(&user)
}

func usersList(c *gin.Context) {
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
