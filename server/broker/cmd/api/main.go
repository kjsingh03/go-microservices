package main

import (
	"context"
	"log"
	"service-broker/internal/app"
)

func main() {

	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	} 

	ctx := context.Background()
	if err := application.Start(ctx); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}