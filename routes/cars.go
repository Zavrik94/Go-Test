package routes

import (
	"github.com/gin-gonic/gin"
	"go-test/database"
	"go-test/database/models"
	"go-test/enums"
	"go-test/middlewares"
	"net/http"
	"time"
)

type CarJSON struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Model string `json:"model"`
}

func RegisterCarsRoutes(router *gin.Engine) {
	router.POST("/cars", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.Admin()), createCar)
	router.POST("/rent/:car_id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.User()), rentCar)
	router.GET("/rented", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.Admin()), rentedCars)
	router.GET("/free", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.User(), enums.Roles{}.Admin()), freeCars)
}

func createCar(c *gin.Context) {
	var car models.Car
	if err := c.BindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()

	if err := db.Create(&car).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, gin.H{})
}

func rentCar(c *gin.Context) {

	carId := c.Param("car_id")

	userID, exist := c.Get("userID")

	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}

	var user models.User

	db := database.GetDB()

	db.Preload("Car").Where("id = ?", userID.(uint)).First(&user)

	if user.Car != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already rent Car"})
		return
	}

	db.Model(&models.Car{}).Where("id = ?", carId).Updates(models.Car{UserID: user.ID})
}

func rentedCars(c *gin.Context) {
	db := database.GetDB()

	var cars []models.Car

	db.Model(&models.Car{}).Where("user_id is not null").Find(&cars)

	var carsJSON []CarJSON
	for _, car := range cars {
		carJSON := CarJSON{
			ID:    car.ID,
			Name:  car.Name,
			Model: car.Model,
		}
		carsJSON = append(carsJSON, carJSON)
	}

	// Return the JSON response
	c.JSON(http.StatusOK, gin.H{
		"cars": carsJSON,
	})
}

func freeCars(c *gin.Context) {
	db := database.GetDB()

	var cars []models.Car

	query := db.Model(&models.Car{}).Where("user_id is null")

	startDateStr := c.Query("start_date")
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			query = query.Where("date >= ?", startDate)
		}
	}

	endDateStr := c.Query("end_date")
	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			query = query.Where("date <= ?", endDate)
		}
	}

	model := c.Query("model")
	if model != "" {
		query = query.Where("model = ?", model)
	}

	manufacturer := c.Query("manufacturer")
	if manufacturer != "" {
		query = query.Where("manufacturer = ?", manufacturer)
	}

	name := c.Query("name")
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Find(&cars)

	var carsJSON []CarJSON
	for _, car := range cars {
		carJSON := CarJSON{
			ID:    car.ID,
			Name:  car.Name,
			Model: car.Model,
		}
		carsJSON = append(carsJSON, carJSON)
	}

	// Return the JSON response
	c.JSON(http.StatusOK, gin.H{
		"cars": carsJSON,
	})
}
