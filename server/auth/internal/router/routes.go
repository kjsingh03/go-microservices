package router

import (
	"authentication/internal/handler"
	"authentication/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		payload := JsonResponse{
			Error:   false,
			Message: "Welcome to Authentication service",
		}

		if err := utils.WriteJSON(w, http.StatusOK, payload); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	})

	mux.Post("/register", handler.RegisterHandler)
	mux.Post("/login", handler.LoginHandler)

	return mux
}
