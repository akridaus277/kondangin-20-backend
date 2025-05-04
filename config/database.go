package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),     // misal: root
		os.Getenv("DB_PASSWORD"), // misal: password
		os.Getenv("DB_HOST"),     // misal: localhost
		os.Getenv("DB_PORT"),     // misal: 3306
		os.Getenv("DB_NAME"),     // misal: gin_hello_world
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	DB = database
}

func GetDB() *gorm.DB {
	return DB
}
