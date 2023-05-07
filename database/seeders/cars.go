package seeders

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/bxcodec/faker/v3"
	"go-test/database/models"
	"math/rand"
	"time"
)

func SeedCars() {
	carModels := map[string][]string{
		"Audi":          {"A1", "A3", "A4", "A6", "Q3"},
		"BMW":           {"1 Series", "3 Series", "5 Series", "7 Series", "X1"},
		"Chevrolet":     {"Camaro", "Corvette", "Impala", "Malibu", "Silverado"},
		"Ford":          {"EcoSport", "Escape", "Expedition", "Explorer", "F-150"},
		"Honda":         {"Accord", "Civic", "CR-V", "Fit", "Odyssey"},
		"Mercedes-Benz": {"A-Class", "C-Class", "E-Class", "S-Class", "GLC"},
		"Nissan":        {"Altima", "Frontier", "Maxima", "Rogue", "Sentra"},
		"Toyota":        {"Camry", "Corolla", "Highlander", "RAV4", "Tacoma"},
		"Volkswagen":    {"Atlas", "Golf", "Jetta", "Passat", "Tiguan"},
		"Volvo":         {"S60", "S90", "V60", "V90", "XC60"},
	}

	for i := 0; i < 10; i++ {
		// select a random brand
		brands := make([]string, 0, len(carModels))
		for brand := range carModels {
			brands = append(brands, brand)
		}
		randomBrand := brands[rand.Intn(len(brands))]

		// select a random model for the selected brand
		randomModel := carModels[randomBrand][rand.Intn(len(carModels[randomBrand]))]

		startDateStr := "2020-01-01"
		endDateStr := "2022-12-31"

		// parse the start and end dates
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		endDate, _ := time.Parse("2006-01-02", endDateStr)

		// generate a random date between the start and end dates
		randomDate := gofakeit.DateRange(startDate, endDate)

		car := models.Car{
			Name:         faker.Word(),
			Model:        randomModel,
			Manufacturer: randomBrand,
			Date:         randomDate,
		}
		car.Create()
	}
}
