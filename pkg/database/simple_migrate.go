package database

import (
	"fmt"
	"gorm.io/gorm"
)

// SimpleAutoMigrate runs GORM's auto-migration for all models
func SimpleAutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&MediaModel{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}
	
	fmt.Println("Database schema migrated successfully")
	return nil
}

// CreateIndexes creates additional indexes that AutoMigrate doesn't handle
func CreateIndexes(db *gorm.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_media_tags ON media_files USING GIN(tags)",
		"CREATE INDEX IF NOT EXISTS idx_media_created_at ON media_files(created_at DESC)",
	}
	
	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}
	
	fmt.Println("Database indexes created successfully")
	return nil
}
