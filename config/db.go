package config

import (
	"MokutonCoffee-API/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Envirnment variable is not set")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = database.AutoMigrate(&models.Product{}, models.User{})
	if err != nil {
		log.Fatal("database migration failed:", err)
	}

	DB = database
	log.Println("Successfully connected to Supabase!")
}
