package models

import (
	"context"
	"database/sql"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

// Create inserts a new user record into the database.
// It returns an error if the insertion fails.
// On success, it updates the user's ID and timestamps fields.
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role, town)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.DB.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Town,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}
