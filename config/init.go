package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_NAME      string
	DB_HOST      string
	DB_PORT      string
	DB_USER      string
	DB_PASSWORD  string
	Port         string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ScopeProfile string
	ScopeEmail   string
	JWTSecret    string
	ServerIp     string
}

var Init Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found!")
	}

	Init = Config{
		DB_NAME:      os.Getenv("DB_NAME"),
		DB_HOST:      os.Getenv("DB_HOST"),
		DB_PORT:      os.Getenv("DB_PORT"),
		DB_USER:      os.Getenv("DB_USER"),
		DB_PASSWORD:  os.Getenv("DB_PASSWORD"),
		Port:         os.Getenv("PORT"),
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		RedirectURL:  os.Getenv("RedirectURL"),
		ScopeProfile: os.Getenv("ScopeProfile"),
		ScopeEmail:   os.Getenv("ScopeEmail"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		ServerIp:     os.Getenv("SERVER_ip"),
	}
}
