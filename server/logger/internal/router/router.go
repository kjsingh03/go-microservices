package router

import (
	"logger/internal/handler"
	"logger/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	logHandler *handlers.LogHandler
}

func NewApp(logHandler *handlers.LogHandler) *App {
	return &App{
		logHandler: logHandler,
	}
}

func (a *App) Routes() http.Handler {
	router := mux.NewRouter()

	// Global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logging())
	router.Use(middleware.Recovery())

	// API v1 routes
	v1 := router.PathPrefix("/api/v1").Subrouter()
	a.setupV1Routes(v1)

	// Health check
	router.HandleFunc("/health", a.logHandler.Health).Methods("GET")
	router.HandleFunc("/", a.logHandler.Home).Methods("GET")

	return router
}

func (a *App) setupV1Routes(router *mux.Router) {
	// Log routes
	logs := router.PathPrefix("/logs").Subrouter()
	logs.HandleFunc("", a.logHandler.GetAllLogs).Methods("GET")
	logs.HandleFunc("", a.logHandler.CreateLog).Methods("POST")
	logs.HandleFunc("/{id}", a.logHandler.GetLogByID).Methods("GET")
	logs.HandleFunc("/{id}", a.logHandler.UpdateLog).Methods("PUT")
	logs.HandleFunc("/{id}", a.logHandler.DeleteLog).Methods("DELETE")
	
	// Utility routes
	logs.HandleFunc("/stats", a.logHandler.GetLogsStats).Methods("GET")
	logs.HandleFunc("/drop", a.logHandler.DropAllLogs).Methods("DELETE").
		Queries("confirm", "true") // Require query parameter for safety
}