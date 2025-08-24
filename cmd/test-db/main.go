package main

import (
	"context"
	"fmt"
	"log"

	"thamaniyah/internal/config"
	"thamaniyah/internal/domain"
	"thamaniyah/internal/repository"
	"thamaniyah/pkg/database"

	"github.com/google/uuid"
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

	// Test database connection
	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Database connection successful!")

	// Create repository
	mediaRepo := repository.NewPostgresMediaRepository(conn)

	// Test creating a media record
	testMedia := &domain.Media{
		ID:          uuid.New().String(),
		Title:       "Test Video",
		Description: "This is a test video",
		FilePath:    "/uploads/test-video.mp4",
		FileSize:    1024000, // 1MB
		Duration:    120,     // 2 minutes
		Format:      "mp4",
		Type:        domain.TypeVideo,
		Status:      domain.StatusUploading,
		Tags:        []string{"test", "demo"},
	}

	ctx := context.Background()

	// Create the media record
	if err := mediaRepo.Create(ctx, testMedia); err != nil {
		log.Fatalf("Failed to create media: %v", err)
	}
	fmt.Printf("Created media with ID: %s\n", testMedia.ID)

	// Get the media record back
	retrievedMedia, err := mediaRepo.GetByID(ctx, testMedia.ID)
	if err != nil {
		log.Fatalf("Failed to get media: %v", err)
	}
	fmt.Printf("Retrieved media: %s - %s\n", retrievedMedia.Title, retrievedMedia.Status)

	// Update status
	if err := mediaRepo.UpdateStatus(ctx, testMedia.ID, domain.StatusReady); err != nil {
		log.Fatalf("Failed to update status: %v", err)
	}
	fmt.Println("Updated media status to ready")

	// Get total count
	total, err := mediaRepo.GetTotal(ctx)
	if err != nil {
		log.Fatalf("Failed to get total: %v", err)
	}
	fmt.Printf("Total media records: %d\n", total)

	fmt.Println("Repository test completed successfully!")
}
