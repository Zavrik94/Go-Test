package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	DeletedAt time.Time `sql:"index"`
}

func (u *User) SoftDelete() bool {
	u.DeletedAt = time.Now()

	return u.Update()
}

func (u *User) Create() bool {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}

	u.Password = string(hashedPassword)

	var count int64
	db.Model(&User{}).Where("email = ?", u.Email).Count(&count)
	if count > 0 {
		return false
	}

	if err := db.Create(&u).Error; err != nil {
		return false
	}

	return true
}

func (u *User) FindByID(id int, relations ...string) {
	query := db.Where("id = ?", id)

	for _, relation := range relations {
		query = query.Preload(relation)
	}

	query.First(&u)
}

func (u *User) Update() bool {
	err := db.Updates(&u).Error

	return err == nil
}

func GetAllUsers() []User {
	var users []User

	db.Preload("Role").Model(&User{}).Find(&users)

	return users
}
