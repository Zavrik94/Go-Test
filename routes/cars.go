package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-test/controllers"
	"go-test/database/models"
	"go-test/enums"
	"go-test/middlewares"
	"net/http"
)

func RegisterCarsRoutes(router *gin.Engine) {
	router.POST("/cars", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.Admin()), validateJSON(), controllers.CreateCar)
	router.POST("/rent/:car_id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.User()), controllers.RentCar)
	router.GET("/rented", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.Admin()), controllers.RentedCars)
	router.GET("/free", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.User(), enums.Roles{}.Admin()), controllers.FreeCars)
}

func validateJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		var car models.Car

		if err := c.ShouldBindJSON(&car); err != nil {

			var verr validator.ValidationErrors
			if errors.As(err, &verr) {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":  "Invalid JSON",
					"errors": verr.Error(),
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("car", car)
		c.Next()
	}
}
