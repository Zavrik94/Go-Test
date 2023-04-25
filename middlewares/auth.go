package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-test/database"
	"go-test/database/models"
	"net/http"
)

// AuthMiddleware Define a middleware function that checks if the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the JWT token is present in the request header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Parse the JWT token and validate the signature
		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key used to sign the token
			return []byte("secret"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check if the token is valid and extract the user ID
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["id"]

			db := database.GetDB()

			var user models.User
			db.Preload("Role").Where("id = ?", userID).First(&user)

			// Add the user ID to the request context
			c.Set("userID", user.ID)

			c.Set("role", user.Role.Name)

			// Call the next handler in the chain
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
