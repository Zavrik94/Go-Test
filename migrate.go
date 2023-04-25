package main

import (
	"fmt"
	"go-test/database"
	"go-test/database/models"
	_ "gorm.io/driver/mysql"
)

func main() {
	db := database.GetDB()

	db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Car{},
	)

	var roles []models.Role

	db.Raw("SELECT * FROM roles").Scan(&roles)

	for _, role := range roles {
		fmt.Printf("%s", role.Name)
	}

	seed()
}

func seed() {
	models.CreateRoles()
}
