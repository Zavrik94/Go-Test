package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID        uint
	Name      string
	Email     string
	Password  string
	RoleID    uint `gorm:"foreignKey:RoleID"`
	Role      Role `gorm:"foreignKey:RoleID"`
	Car       *Car
	DeletedAt *time.Time `sql:"index"`
}

func (User) SoftDelete() bool {
	return true
}
