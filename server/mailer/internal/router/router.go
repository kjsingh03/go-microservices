package router

import (
	"mailer/internal/handler"
	"mailer/internal/middleware"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Routes(h *handler.Handler) http.Handler {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	
	router.Use(middleware.RecoveryMiddleware)

	router.Use(middleware.RequestIDMiddleware)

	api := router.PathPrefix("/api/v1").Subrouter()
	
	router.HandleFunc("/", h.Home).Methods("GET")
	
	api.Use(middleware.RateLimitMiddleware(100, 60)) // 100 requests per minute
	api.HandleFunc("/send", h.SendMail).Methods("POST")
	
	api.HandleFunc("/send/batch", h.SendBatchMail).Methods("POST")

	return setupCORS(router)
}

func setupCORS(router *mux.Router) http.Handler {
	routerConfig := DefaultRouterConfig()

	if os.Getenv("ENVIRONMENT") == "development" {
		routerConfig.Debug = true
		routerConfig.AllowedOrigins = []string{"*"}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   routerConfig.AllowedOrigins,
		AllowedMethods:   routerConfig.AllowedMethods,
		AllowedHeaders:   routerConfig.AllowedHeaders,
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300, 
		Debug:            routerConfig.Debug,
	})

	return c.Handler(router)
}

type RouterConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	Debug          bool
}

func DefaultRouterConfig() *RouterConfig {
	return &RouterConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001", 
			"http://localhost:5173",
		},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"X-CSRF-Token",
		},
		Debug: false,
	}
}