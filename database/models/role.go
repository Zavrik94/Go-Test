package models

import (
	"go-test/database"
	"go-test/enums"
)

type Role struct {
	ID   uint
	Name string
}

func CreateRoles() {
	db := database.GetDB()

	roles := []Role{
		{Name: enums.Roles{}.User()},
		{Name: enums.Roles{}.Admin()},
		{Name: enums.Roles{}.SuperAdmin()},
	}

	for _, role := range roles {
		db.FirstOrCreate(&role, Role{Name: role.Name})
	}
}
