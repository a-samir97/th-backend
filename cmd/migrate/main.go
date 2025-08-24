package main

import (
	"fmt"
	"log"

	"thamaniyah/internal/config"
	"thamaniyah/pkg/database"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	conn, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Run simple auto-migration
	fmt.Println("Running GORM auto-migration...")
	if err := database.SimpleAutoMigrate(conn.DB); err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}

	// Create additional indexes
	fmt.Println("Creating additional indexes...")
	if err := database.CreateIndexes(conn.DB); err != nil {
		log.Fatalf("Index creation failed: %v", err)
	}

	fmt.Println("Database setup completed successfully!")
}
