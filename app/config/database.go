package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const dsn = "root:password@tcp(db:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"

func NewDBConn() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
