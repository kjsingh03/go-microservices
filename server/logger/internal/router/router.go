package router

import (
	"logger/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Routes() http.Handler {
	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Health check route
	router.HandleFunc("/", handler.Home).Methods("GET")

	// Log management routes
	router.HandleFunc("/logs", handler.GetAllLogs).Methods("GET")
	router.HandleFunc("/logs", handler.CreateLog).Methods("POST")
	router.HandleFunc("/logs/{id}", handler.GetLogByID).Methods("GET")
	router.HandleFunc("/logs/{id}", handler.UpdateLog).Methods("PUT")
	router.HandleFunc("/logs/{id}", handler.DeleteLog).Methods("DELETE")
	
	// Additional utility routes
	router.HandleFunc("/logs/stats", handler.GetLogsStats).Methods("GET")
	router.HandleFunc("/logs/drop", handler.DropAllLogs).Methods("DELETE")

	return c.Handler(router)
}