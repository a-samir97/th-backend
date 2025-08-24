package database

import (
	"fmt"
	"log"

	"thamaniyah/internal/domain"

	"gorm.io/gorm"
)

// SimpleAutoMigrate performs simple auto-migration for domain models
func SimpleAutoMigrate(db *gorm.DB) error {
	// Auto-migrate all domain models
	err := db.AutoMigrate(
		&domain.Media{},
		&domain.SearchIndex{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	log.Println("✓ Database migration completed successfully")
	return nil
}

// CreateIndexes creates additional indexes that AutoMigrate doesn't handle
func CreateIndexes(db *gorm.DB) error {
	// Drop existing indexes that might conflict
	dropIndexes := []string{
		"DROP INDEX IF EXISTS idx_media_tags",
	}

	for _, dropSQL := range dropIndexes {
		if err := db.Exec(dropSQL).Error; err != nil {
			log.Printf("Warning: Failed to drop index: %v", err)
		}
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_media_created_at ON media_files(created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_search_content ON search_index USING GIN(to_tsvector('english', content))",
		"CREATE INDEX IF NOT EXISTS idx_search_media_id ON search_index(media_id)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	fmt.Println("Database indexes created successfully")
	return nil
}
