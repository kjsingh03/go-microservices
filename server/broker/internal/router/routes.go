// internal/router/router.go
package router

import (
	"net/http"
	"service-broker/internal/handler"
	"service-broker/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func New(h *handler.Handler) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.CorsMiddleware())

	mux.Get("/", h.Home)
	mux.Post("/handle", h.HandleSubmission)

	return mux
}