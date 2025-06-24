package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()
		
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(wrapped, r)
		
		// duration := time.Since(start)
		// requestID := r.Header.Get("X-Request-ID")
		
		// log.Printf("[%s] %s %s %d %v %s",
		// 	requestID,
		// 	r.Method,
		// 	r.URL.Path,
		// 	wrapped.statusCode,
		// 	duration,
		// 	r.RemoteAddr,
		// )
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Add to response header
		w.Header().Set("X-Request-ID", requestID)
		
		// Add to request header for downstream handlers
		r.Header.Set("X-Request-ID", requestID)
		
		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(requests int, duration int) mux.MiddlewareFunc {
	limiter := rate.NewLimiter(rate.Every(time.Duration(duration)*time.Second/time.Duration(requests)), requests)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}