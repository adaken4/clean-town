package database

import (
	"context"
	"database/sql"
	"time"
)

func CheckDB(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}
