package config

import (
	"log"
	"os"

	"example.com/event-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URI")

	if dsn == "" {
		log.Fatal("Environment Variable Belum Di isi")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal Terkoneksi Di Database", err)
	}

	err = database.AutoMigrate(&models.Event{}, &models.User{}, &models.Booking{})

	if err != nil {
		log.Fatal("Gagal Melakukan Migration Ke Database :", err)
	}

	DB = database
	log.Println("Berhasil Terkoneksi Ke Database")
}
