package database

import (
	"fmt"
	"time"

	"thamaniyah/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connection holds the database connection
type Connection struct {
	DB *gorm.DB
}

// NewPostgresConnection creates a new PostgreSQL connection
func NewPostgresConnection(cfg *config.Config) (*Connection, error) {
	dsn := cfg.DatabaseURL()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Connection{DB: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping checks if the database is reachable
func (c *Connection) Ping() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Health returns the database health status
func (c *Connection) Health() map[string]interface{} {
	sqlDB, _ := c.DB.DB()
	stats := sqlDB.Stats()

	return map[string]interface{}{
		"status":             "up",
		"open_connections":   stats.OpenConnections,
		"idle_connections":   stats.Idle,
		"in_use_connections": stats.InUse,
	}
}

// Transaction executes a function within a database transaction
func (c *Connection) Transaction(fn func(*gorm.DB) error) error {
	return c.DB.Transaction(fn)
}
