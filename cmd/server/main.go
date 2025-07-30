package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"golang-woktopup/config"
	"golang-woktopup/internal/model"
	"golang-woktopup/internal/router"
)

func main() {
	// 1. Load .env terlebih dahulu
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Baru inisialisasi Midtrans setelah env dimuat
	config.InitMidtrans()

	// 3. Koneksi DB dan lainnya
	db := config.ConnectDB()
	db.AutoMigrate(&model.User{})
	r := router.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
