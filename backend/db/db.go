package db

import (
	"log"
	"os"

	"github.com/Similadayo/backend/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func ConnectToDb() {

	Err := godotenv.Load()
	if Err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_Name")

	connectionString := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	// Automatically create tables
	DB.AutoMigrate(&models.Order{}, models.Product{}, models.User{}, models.Cart{}, models.Category{}, models.Role{})
}
