package seeders

import (
	"github.com/bxcodec/faker/v3"
	"go-test/database/models"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

func SeedUsers() {
	for i := 0; i < 5; i++ {
		Password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

		user := models.User{
			Email:    faker.Email(),
			Name:     faker.Name(),
			RoleID:   uint(rand.Intn(2)) + 1,
			Password: string(Password),
		}
		user.Create()
	}
}
