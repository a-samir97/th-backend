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
	"thamaniyah/pkg/httpclient"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database (same database, different service)
	conn, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Initialize HTTP client for CMS service communication
	cmsClient := httpclient.NewClient(fmt.Sprintf("http://localhost:%d", cfg.Server.Port))

	// Initialize repositories
	searchRepo := repository.NewPostgresSearchRepository(conn)

	// Initialize services
	searchService := service.NewSearchService(searchRepo, cmsClient)

	// Initialize handlers
	searchHandler := handler.NewSearchHandler(searchService)

	// Setup router
	router := setupRouter(searchHandler)

	// Start server on different port (8081)
	discoveryPort := cfg.Server.Port + 1
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", discoveryPort),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Discovery Service starting on port %d", discoveryPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Discovery service...")

	// Graceful shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Discovery Service shutdown complete")
}

// setupRouter configures the HTTP router with routes and middleware
func setupRouter(searchHandler *handler.SearchHandler) *gin.Engine {
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
			"service":   "discovery-service",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		search := v1.Group("/search")
		{
			search.GET("", searchHandler.Search)
			search.GET("/suggest", searchHandler.Suggest)
			search.POST("/reindex", searchHandler.Reindex)
		}
	}

	return router
}
