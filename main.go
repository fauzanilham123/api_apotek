package main

import (
	"api_apotek/config"
	"api_apotek/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// for load godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// database connection
	db := config.ConnectDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// router
	r := routes.SetupRouter(db)
	r.Run(":8000")
}