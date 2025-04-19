package config

import (
	"log"
	"main/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var DB *gorm.DB

var host, user, password, dbName, port, sslMode string

func initENVVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
	user = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	sslMode = os.Getenv("DB_SSL_MODE")
}

func InitDB() {
	initENVVariables()

	dsn := "host=" + host +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbName +
		" port=" + port +
		" sslmode=" + sslMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.Exec("CREATE SCHEMA IF NOT EXISTS public")

	db.AutoMigrate(&models.Customers{}, &models.Invoices{})

	DB = db
}
