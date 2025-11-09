package models

import (
	"context"
	"time"
)

// UserRepository defines the interface for user-related database operations.
// Implementations should ensure data integrity including handling soft deletes if applicable.
type UserRepository interface {
	Create(ctx context.Context, user *User, token string, expiry time.Time) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error // For profile completion
	UpdateStatus(ctx context.Context, userID uint64, status string) error

	FindByTown(ctx context.Context, town string, role string) ([]*User, error)
	FindBySkills(ctx context.Context, skills []string) ([]*User, error)
	MarkUserAsVerified(ctx context.Context, userID uint64) error
}
