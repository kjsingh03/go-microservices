// internal/app/app.go
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service-broker/internal/config"
	"service-broker/internal/handler"
	"service-broker/internal/router"
	"service-broker/internal/service"
	"syscall"
	"time"
)

type App struct {
	Config   *config.Config
	Services *service.Services
	Server   *http.Server
}

func New() (*App, error) {

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	services, err := initServices(cfg)
	if err != nil {
		return nil, err
	}

	handlers := handler.New(services)

	r := router.New(handlers)

	server := &http.Server{
		Addr:    ":"+cfg.Server.Port,
		Handler: r,
	}

	return &App{
		Config:   cfg,
		Services: services,
		Server:   server,
	}, nil
}

func initServices(cfg *config.Config) (*service.Services, error) {
	rabbitService, err := service.NewRabbitService(cfg.Rabbit)

	if err != nil {
		return nil, err
	}

	authService := service.NewAuthService(cfg.Services.AuthURL, cfg.Services.Timeout, cfg.Services.RetryCount, "")
	logService := service.NewLogService(cfg.Services.LogURL, cfg.Services.Timeout, cfg.Services.RetryCount)
	mailService := service.NewMailService(cfg.Services.MailURL, cfg.Services.Timeout, cfg.Services.RetryCount, "")

	return &service.Services{
		AuthService:   authService,
		LogService:    logService,
		MailService:   mailService,
		RabbitService: rabbitService,
	}, nil
}

func (a *App) Start(ctx context.Context) error {

	go func() {
		log.Printf("Server starting at http://localhost:%s", a.Config.Server.Port)
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down server...")
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	if err := a.Services.Close(); err != nil {
		log.Printf("Services close error: %v", err)
	}

	log.Println("Server stopped")
	return nil
}
