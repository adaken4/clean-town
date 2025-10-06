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

// FindByEmail retrieves a user by their email address (case-insensitive via CITEXT).
// It excludes soft-deleted users by ensuring deleted_at IS NULL.
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	user := new(User)
	query := `
        SELECT id, name, email, password_hash, role, town, status, created_at, updated_at, deleted_at
        FROM users
        WHERE email = $1 AND deleted_at IS NULL
    `
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.Town, &user.Status,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByID retrieves a user by their unique ID.
// It excludes soft-deleted users by filtering on deleted_at IS NULL.
func (r *PostgresUserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	user := new(User)
	query := `
        SELECT id, name, email, password_hash, role, town, status, created_at, updated_at, deleted_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.Town, &user.Status,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
