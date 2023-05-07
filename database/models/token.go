package models

import (
	"go-test/database"
	"gorm.io/gorm/clause"
	"time"
)

type Token struct {
	ID        uint
	ExpiredAt time.Time
	UserID    *uint `gorm:"ForeignKey:user_id;uniqueIndex"`
}

func init() {
	db = database.GetDB()
}

func (t *Token) Save() bool {
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"expired_at"}),
	}).Create(&t).Error

	return err == nil
}

func DeleteToken(UserId uint) bool {
	err := db.Where("user_id = ?", UserId).Delete(&Token{})

	return err == nil
}

func CheckToken(UserId uint) bool {
	var token Token
	res := db.Where("user_id = ?", UserId).First(&token)

	if res.Error != nil || res.RowsAffected == 0 {
		return false
	}

	if token.ExpiredAt.Before(time.Now()) {
		return false
	}

	return true
}
