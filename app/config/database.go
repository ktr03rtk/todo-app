package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBConn() *gorm.DB {
	dsn, err := getDsn()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func getDsn() (string, error) {
	dbUser, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		return "", errors.New("env DB_USERNAME is not found")
	}

	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return "", errors.New("env DB_PASSWORD is not found")
	}

	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return "", errors.New("env DB_HOST is not found")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return "", errors.New("env DB_NAME is not found")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName)

	return dsn, nil
}
