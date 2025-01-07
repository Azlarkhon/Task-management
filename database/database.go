package database

import (
	"fmt"
	"log"
	"task-management/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		config.Init.DB_HOST, config.Init.DB_PORT, config.Init.DB_USER, config.Init.DB_NAME, config.Init.DB_PASSWORD)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	log.Println("Database connection successfully initialized")
}
