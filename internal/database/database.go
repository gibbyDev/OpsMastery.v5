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
// This variable holds the database connection and is shared across the app.
var db *gorm.DB

// DB returns the global db instance
// Use this function to get the database connection anywhere in your codebase.
func DB() *gorm.DB {
	return db
}

// Init initializes the global db instance and runs automigrate
// This function sets up the database connection using environment variables,
// and automatically migrates all models to ensure the schema is up to date.
func Init() {
	// Load environment variables for database connection.
	// These should be set in your .env file or container environment.
	dbHost := os.Getenv("BLUEPRINT_DB_HOST")         // Database host (e.g., psql_bp or localhost)
	dbUser := os.Getenv("BLUEPRINT_DB_USERNAME")     // Database username
	dbPassword := os.Getenv("BLUEPRINT_DB_PASSWORD") // Database password
	dbName := os.Getenv("BLUEPRINT_DB_DATABASE")     // Database name
	dbPort := os.Getenv("BLUEPRINT_DB_PORT")         // Database port (usually 5432 for Postgres)
	dbSSLMode := "disable"                           // SSL mode for local/dev; can be set via env for prod
	dbSchema := os.Getenv("BLUEPRINT_DB_SCHEMA")     // Database schema (e.g., public)

	// Build the Data Source Name (DSN) string for Postgres connection.
	// This string contains all connection parameters.
	dsn := "host=" + dbHost +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" port=" + dbPort +
		" sslmode=" + dbSSLMode +
		" search_path=" + dbSchema

	// Attempt to open a connection to the database using GORM and the DSN.
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// If connection fails, log the error and exit the application.
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Automigrate all models.
	// This will create or update tables for User, Ticket, Client, Chat, and ChatMessage.
	// It ensures your database schema matches your Go models.
	if err := db.AutoMigrate(
		&models.User{},
		&models.Ticket{},
		&models.Client{},
		&models.Chat{},
		&models.ChatMessage{},
	); err != nil {
		// If migration fails, log the error and exit the application.
		log.Fatalf("failed to auto migrate models: %v", err)
	}
}
