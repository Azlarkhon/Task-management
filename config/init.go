package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_NAME     string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	Port        string
}

var Init Config
var JWTSecret string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found!")
	}

	Init = Config{
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		Port:        os.Getenv("PORT"),
	}

	JWTSecret = os.Getenv("JWT_SECRET")
}
