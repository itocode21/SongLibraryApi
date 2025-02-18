package repositories

import (
	"SongLibraryApi/models"
	"SongLibraryApi/utils"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	config := utils.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

	var err error
	log.Printf("Connecting to DB with DSN: %s", dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	err = DB.AutoMigrate(&models.Song{})
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Connection to database successfully")
}
