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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.ConnectDB()
	db.AutoMigrate(&model.User{})
	r := router.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
