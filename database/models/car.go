package models

import (
	"go-test/database"
	"time"
)

type Car struct {
	ID           uint
	Name         string    `json:"name" binding:"required"`
	Model        string    `json:"model" binding:"required"`
	Manufacturer string    `json:"manufacturer" binding:"required"`
	Date         time.Time `gorm:"type:date" json:"date" binding:"required" time_format:"2006-01-02"`
	UserID       *uint     `gorm:"null;ForeignKey:user_id"`
}

func init() {
	db = database.GetDB()
}

func (c *Car) Create() bool {
	c.UserID = nil
	err := db.Create(&c).Error

	return err == nil
}

func (c *Car) Update() bool {
	err := db.Updates(&c).Error

	return err == nil
}

func (c *Car) Delete() bool {
	err := db.Delete(&c).Error

	return err == nil
}
