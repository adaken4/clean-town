package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectWithRetry(dsn string, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error
	backoff := time.Second

	// Exponential backoff retry
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("pgx", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			// Configure pool
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(25)
			db.SetConnMaxLifetime(5 * time.Minute)

			return db, nil
		}

		log.Printf("db connect attempt %d failed: %v", i+1, err)
		time.Sleep(backoff)
		backoff *= 2 // exponential growth
	}

	return nil, fmt.Errorf("could not connect to db after %d attempts: %w", maxRetries, err)
}
