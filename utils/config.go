package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBUser      string
	DBPassword  string
	DBName      string
	DBHost      string
	DBPort      string
	ExternalAPI string
	LogLevel    string
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		ExternalAPI: os.Getenv("EXTERNAL_API_URL"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
	}
}
