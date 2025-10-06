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

// Update modifies the profile fields of an existing user.
// This excludes password changes and soft-deleted users.
// The updated_at timestamp is refreshed.
func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	query := `
        UPDATE users
        SET name = $1, email = $2, role = $3, town = $4, status = $5, updated_at = NOW()
        WHERE id = $6 AND deleted_at IS NULL
    `
	_, err := r.DB.ExecContext(ctx, query, user.Name, user.Email, user.Role, user.Town, user.Status, user.ID)
	return err
}

// UpdateStatus updates only the status field of a user,
// ensuring the user is not soft-deleted. It refreshes updated_at.
func (r *PostgresUserRepository) UpdateStatus(ctx context.Context, userID uint64, status string) error {
	query := `
        UPDATE users
        SET status = $1, updated_at = NOW()
        WHERE id = $2 AND deleted_at IS NULL
    `
	_, err := r.DB.ExecContext(ctx, query, status, userID)
	return err
}

// FindByTown fetches all users located in a specific town with a given role,
// excluding soft-deleted records.
func (r *PostgresUserRepository) FindByTown(ctx context.Context, town string, role string) ([]*User, error) {
	query := `
        SELECT id, name, email, password_hash, role, town, status, created_at, updated_at, deleted_at
        FROM users
        WHERE town = $1 AND role = $2 AND deleted_at IS NULL
    `
	rows, err := r.DB.QueryContext(ctx, query, town, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := new(User)
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Role, &user.Town, &user.Status,
			&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
