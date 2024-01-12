package config

import (
	"api_apotek/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
    // koneksi dari ENV
    username := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    database := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%v:%v@%v/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, database)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
        panic(err.Error())
    }

	db.AutoMigrate(&models.User{}, &models.LogActivity{}, &models.Drug{}, &models.Transaction{}, &models.Recipe{})

    
    return db
}