package database

import (
	"log"
	"os"

	"OpsMastery.v5/internal/models"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global GORM DB instance
var db *gorm.DB

// DB returns the global db instance
func DB() *gorm.DB {
	return db
}

// Init initializes the global db instance and runs automigrate
func Init() {
	// Load environment variables
	dbHost := os.Getenv("BLUEPRINT_DB_HOST")
	dbUser := os.Getenv("BLUEPRINT_DB_USERNAME")
	dbPassword := os.Getenv("BLUEPRINT_DB_PASSWORD")
	dbName := os.Getenv("BLUEPRINT_DB_DATABASE")
	dbPort := os.Getenv("BLUEPRINT_DB_PORT")
	dbSSLMode := "disable" // or os.Getenv("BLUEPRINT_DB_SSLMODE")
	dbSchema := os.Getenv("BLUEPRINT_DB_SCHEMA")

	dsn := "host=" + dbHost +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" port=" + dbPort +
		" sslmode=" + dbSSLMode +
		" search_path=" + dbSchema

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Automigrate all models
	if err := db.AutoMigrate(&models.User{}, &models.Ticket{}, &models.Client{}, &models.Chat{}, &models.ChatMessage{}); err != nil {
		log.Fatalf("failed to auto migrate models: %v", err)
	}
}
