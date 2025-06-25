package main

import (
	"authentication/internal/config"
	"authentication/internal/db"
	"authentication/internal/model"
	"authentication/internal/router"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	config.LoadEnv()

	dbpool, err := db.ConnectToDB()

	AUTH_PORT := config.GetPort("AUTH_PORT", "80")

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer dbpool.Close()

	model.SetDB(dbpool)

	log.Printf("Auth Service starting at: http://localhost:%s", AUTH_PORT)

	srv := &http.Server{
		Addr:    ":" + AUTH_PORT,
		Handler: router.Routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
