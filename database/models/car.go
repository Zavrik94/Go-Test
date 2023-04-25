package models

import "time"

type Car struct {
	ID           uint
	Name         string
	Model        string
	Manufacturer string
	Date         time.Time `gorm:"type:date"`
	UserID       uint      `gorm:"ForeignKey:user_id"`
}
