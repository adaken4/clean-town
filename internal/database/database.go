package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/adaken4/clean-town/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectWithRetry(cfg config.DatabaseConfig) (*sql.DB, error) {
	var db *sql.DB
	var openErr, pingErr error
	backoff := time.Duration(cfg.RetryBackoff) * time.Second

	// Exponential backoff retry
	for i := 0; i < cfg.MaxRetries; i++ {
		db, openErr = sql.Open("pgx", cfg.DSN())
		if openErr != nil {
			log.Printf("db open attempt %d failed: %v", i+1, openErr)
		} else {
			pingErr = CheckDB(db)
			if pingErr == nil {
				// Configure pool
				db.SetMaxOpenConns(cfg.MaxOpenConns)
				db.SetMaxIdleConns(cfg.MaxIdleConns)
				db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Minute)

				return db, nil
			}
			db.Close()
			log.Printf("db ping attempt %d failed: %v", i+1, pingErr)
		}

		time.Sleep(backoff)
		backoff *= 2 // exponential growth

		const maxBackoff = 30 * time.Second
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
	// Return the most recent error
	if pingErr != nil {
		return nil, fmt.Errorf("could not connect to db after %d attempts (last ping error): %w", cfg.MaxRetries, pingErr)
	}

	return nil, fmt.Errorf("could not connect to db after %d attempts (last open error): %w", cfg.MaxRetries, openErr)
}
