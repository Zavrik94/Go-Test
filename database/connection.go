package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func init() {
	once.Do(func() {
		dsn := "root:root@tcp(go_mysql:3306)/example_db?charset=utf8mb4&parseTime=True&loc=Local"
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})
}

func GetDB() *gorm.DB {
	return db
}
