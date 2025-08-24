package main

import (
	"fmt"
	"log"
	"os"

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

	// Create migrator
	migrator := database.NewMigrator(conn.DB, "./migrations")

	// Check command line arguments
	if len(os.Args) > 1 && os.Args[1] == "auto" {
		// Run GORM auto-migration
		fmt.Println("Running GORM auto-migration...")
		if err := migrator.AutoMigrate(); err != nil {
			log.Fatalf("Auto-migration failed: %v", err)
		}
		fmt.Println("Auto-migration completed successfully")
		return
	}

	// Run SQL migrations
	fmt.Println("Running SQL migrations...")
	if err := migrator.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Database migration completed successfully!")
}
