package controllers

import (
	"github.com/gin-gonic/gin"
	"go-test/database/models"
	"net/http"
	"strconv"
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

	carId, _ := strconv.Atoi(c.Param("car_id"))

	userID, exist := c.Get("userID")

	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}

	var user models.User

	var car models.Car

	user.FindByID(userID.(int), "Car")

	if user.Car != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already rent Car"})
		return
	}

	car.FindByID(carId)

	car.UserID = &user.ID

	car.Update()
}

type CarJSON struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Model string `json:"model"`
}

func RentedCars(c *gin.Context) {

	conditions := []string{
		"date >= ?",
		"make = ?",
		"year >= ?",
	}

	parameters := []interface{}{
		"not null",
	}

	cars := models.FindAllCars(conditions, parameters)

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
	var conditions []string
	var parameters []interface{}

	startDateStr := c.Query("start_date")
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			conditions = append(conditions, "date >= ?")
			parameters = append(parameters, startDate)
		}
	}

	endDateStr := c.Query("end_date")
	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			conditions = append(conditions, "date <= ?")
			parameters = append(parameters, endDate)
		}
	}

	model := c.Query("model")
	if model != "" {
		conditions = append(conditions, "model = ?")
		parameters = append(parameters, model)
	}

	manufacturer := c.Query("manufacturer")
	if manufacturer != "" {
		conditions = append(conditions, "manufacturer = ?")
		parameters = append(parameters, manufacturer)
	}

	name := c.Query("name")
	if name != "" {
		conditions = append(conditions, "name LIKE ?")
		parameters = append(parameters, "%"+name+"%")
	}

	conditions = append(conditions, "user_id is ?")
	parameters = append(parameters, "null")

	cars := models.FindAllCars(conditions, parameters)

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
