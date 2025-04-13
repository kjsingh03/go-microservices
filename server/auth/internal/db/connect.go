package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"time"
)

func ConnectToDB() (*sql.DB, error) {

	dsn := os.Getenv("POSTGRES_URL")

	maxRetries := 5
	retryInterval := 2 * time.Second

	var db *sql.DB
	var err error

	for retries := 0; retries < maxRetries; retries++ {

		db, err = sql.Open("postgres", dsn)

		if err != nil {
			time.Sleep(retryInterval)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			time.Sleep(retryInterval)
			continue
		}

		log.Println("Connected to Postgres successfully")
		return db, nil
	}

	return nil, fmt.Errorf("could not connect to database after %d retries: %w", maxRetries, err)
}