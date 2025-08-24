package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// Migrator handles database migrations
type Migrator struct {
	db            *gorm.DB
	migrationsDir string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Migration represents a database migration
type Migration struct {
	Version string
	Name    string
	SQL     string
}

// RunMigrations executes all pending migrations
func (m *Migrator) RunMigrations() error {
	// Create migrations tracking table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Filter out already applied migrations
	pendingMigrations := m.filterPendingMigrations(migrations, appliedMigrations)

	// Execute pending migrations
	for _, migration := range pendingMigrations {
		if err := m.executeMigration(migration); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
		}
		fmt.Printf("Applied migration: %s - %s\n", migration.Version, migration.Name)
	}

	if len(pendingMigrations) == 0 {
		fmt.Println("No pending migrations to apply")
	} else {
		fmt.Printf("Applied %d migrations successfully\n", len(pendingMigrations))
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`
	return m.db.Exec(sql).Error
}

// loadMigrations loads all migration files from the migrations directory
func (m *Migrator) loadMigrations() ([]Migration, error) {
	files, err := os.ReadDir(m.migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Parse filename: 001_create_media_files_table.sql
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}

		version := parts[0]
		name := strings.TrimSuffix(parts[1], ".sql")

		// Read file content
		content, err := os.ReadFile(filepath.Join(m.migrationsDir, file.Name()))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			SQL:     string(content),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations returns a list of already applied migration versions
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var versions []string
	err := m.db.Table("schema_migrations").Select("version").Find(&versions).Error
	if err != nil {
		return nil, err
	}

	applied := make(map[string]bool)
	for _, version := range versions {
		applied[version] = true
	}

	return applied, nil
}

// filterPendingMigrations filters out already applied migrations
func (m *Migrator) filterPendingMigrations(migrations []Migration, applied map[string]bool) []Migration {
	var pending []Migration
	for _, migration := range migrations {
		if !applied[migration.Version] {
			pending = append(pending, migration)
		}
	}
	return pending
}

// executeMigration executes a single migration
func (m *Migrator) executeMigration(migration Migration) error {
	// Execute the migration in a transaction
	return m.db.Transaction(func(tx *gorm.DB) error {
		// Execute the migration SQL
		if err := tx.Exec(migration.SQL).Error; err != nil {
			return err
		}

		// Record the migration as applied
		return tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", migration.Version).Error
	})
}

// AutoMigrate runs GORM auto-migration for models
func (m *Migrator) AutoMigrate() error {
	return m.db.AutoMigrate(&MediaModel{})
}
