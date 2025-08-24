package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"thamaniyah/internal/config"
	"thamaniyah/internal/handler"
	"thamaniyah/internal/middleware"
	"thamaniyah/internal/repository"
	"thamaniyah/internal/service"
	"thamaniyah/pkg/database"

	"github.com/gin-gonic/gin"
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

	// Auto-migrate database (for development)
	if err := database.SimpleAutoMigrate(conn.DB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	if err := database.CreateIndexes(conn.DB); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Initialize repositories
	mediaRepo := repository.NewPostgresMediaRepository(conn)

	// Initialize services
	mediaService := service.NewMediaService(mediaRepo)

	// Initialize handlers
	mediaHandler := handler.NewMediaHandler(mediaService)

	// Setup router
	router := setupRouter(mediaHandler)

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("CMS Service starting on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("CMS Service shutdown complete")
}

// setupRouter configures the HTTP router with routes and middleware
func setupRouter(mediaHandler *handler.MediaHandler) *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "cms-service",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		media := v1.Group("/media")
		{
			media.POST("/upload-url", mediaHandler.CreateUploadURL)
			media.POST("/:id/confirm", mediaHandler.ConfirmUpload)
			media.GET("", mediaHandler.GetAllMedia)
			media.GET("/:id", mediaHandler.GetMedia)
			media.PUT("/:id", mediaHandler.UpdateMedia)
			media.DELETE("/:id", mediaHandler.DeleteMedia)
		}
	}

	return router
}
