package main

import (
	"github.com/joho/godotenv"
	"go-test/database/models"
	"go-test/database/seeders"
	_ "gorm.io/driver/mysql"
)

func main() {
	godotenv.Load()

	//db := database.GetDB()
	//
	//db.AutoMigrate(
	//	&models.User{},
	//	&models.Role{},
	//	&models.Car{},
	//	&models.Token{},
	//)
	//
	//seed()
	//
	//if len(os.Args) > 1 {
	//	if os.Args[1] == "--seed" {
	//		seedTestData()
	//	}
	//}

}

func seed() {
	models.CreateRoles()
}

func seedTestData() {
	seeders.SeedUsers()
	seeders.SeedCars()
}
