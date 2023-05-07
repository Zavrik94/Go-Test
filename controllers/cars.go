package controllers

import (
	"github.com/gin-gonic/gin"
	"go-test/database"
	"go-test/database/models"
	"net/http"
	"time"
)

func CreateCar(c *gin.Context) {
	car, exist := c.Get("car")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}

	if carObj, ok := car.(models.Car); ok {
		if carObj.Create() {
			c.JSON(http.StatusCreated, gin.H{})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
	return
}

func RentCar(c *gin.Context) {

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

	db.Model(&models.Car{}).Where("id = ?", carId).Updates(models.Car{UserID: &user.ID})
}

type CarJSON struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Model string `json:"model"`
}

func RentedCars(c *gin.Context) {
	db := database.GetDB()

	var cars []models.Car

	db.Model(&models.Car{}).Where("user_id is not null").Find(&cars)

	carsJSON := make([]CarJSON, len(cars))

	for i, car := range cars {
		carsJSON[i] = CarJSON{
			ID:    car.ID,
			Name:  car.Name,
			Model: car.Model,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"cars": carsJSON,
	})
}

func FreeCars(c *gin.Context) {
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

	carsJSON := make([]CarJSON, len(cars))
	for i, car := range cars {
		carsJSON[i] = CarJSON{
			ID:    car.ID,
			Name:  car.Name,
			Model: car.Model,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"cars": carsJSON,
	})
}
