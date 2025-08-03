package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"golang-woktopup/config"
	"golang-woktopup/internal/model"
	"golang-woktopup/internal/router"
	"golang-woktopup/internal/seed"
)

func main() {
	// 1. Load konfigurasi dari file .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Inisialisasi Midtrans dengan Server Key dari .env
	config.InitMidtrans()

	// 3. Koneksi ke database
	db := config.ConnectDB()

	// 4. AutoMigrate semua tabel
	if err := db.AutoMigrate(
		&model.User{},
		&model.Game{},
		&model.Product{},
		&model.Voucher{},
		&model.Order{},
		&model.Payment{},
		&model.Invoice{},
	); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	// 5. Seed data awal (jika database kosong)
	seed.SeedGamesAndProducts(db)

	// 6. Jalankan server
	r := router.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s", port)
	r.Run(":" + port)
}
