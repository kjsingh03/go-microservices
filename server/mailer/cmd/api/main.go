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

	"mailer/internal/config"
	"mailer/internal/handler"
	"mailer/internal/mailer"
	"mailer/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	mailerService := mailer.NewService(&cfg.Mailer)
	h := handler.NewHandler(mailerService)

	routes := router.Routes(h)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      routes,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("Mailer Service starting at http://localhost:%s", cfg.Server.Port)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	} else{
		log.Println("Server exited gracefully")
	}
}