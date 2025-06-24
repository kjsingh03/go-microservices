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

	"logger/internal/config"
	"logger/internal/db"
	"logger/internal/handler"
	"logger/internal/repositories"
	"logger/internal/router"
	"logger/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dbManager, err := database.NewManager(cfg.Database.URL, cfg.Database.Name)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbManager.Close()

	// Initialize dependencies
	app, err := initializeApp(dbManager)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start server
	if err := startServer(cfg, app); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initializeApp(dbManager *database.Manager) (*router.App, error) {
	// Initialize repository
	logRepo := repositories.NewLogRepository(dbManager.GetDatabase())

	// Initialize service
	logService := services.NewLogService(logRepo)

	// Initialize handlers
	logHandler := handlers.NewLogHandler(logService) // Fixed package name

	// Create app with all handlers
	return router.NewApp(logHandler), nil
}

func startServer(cfg *config.Config, app *router.App) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      app.Routes(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Logger Service starting at: http://localhost:%s", cfg.Server.Port)
		log.Println("Connected to MongoDB successfully")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server exited")
	return nil
}