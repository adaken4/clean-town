package models

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID           int64   `db:"id" json:"id"`
	Name         string  `db:"name" json:"name" validate:"required"`
	Email        string  `db:"email" json:"email" validate:"required,email"`
	Phone        *string `db:"phone" json:"phone" validate:"omitempty,e164"`
	Password     string  `db:"-" json:"password,omitempty" validate:"required,min=8"` // only for input
	PasswordHash []byte  `db:"password_hash" json:"-"`

	// Profile fields
	Town      *string        `db:"town" json:"town"`
	Bio       *string        `db:"bio" json:"bio"`
	Age       *int           `db:"age" json:"age" validate:"omitempty,min=1,max=120"`
	AvatarURL *string        `db:"avatar_url" json:"avatar_url,omitempty"`
	Skills    pq.StringArray `db:"skills" json:"skills"` // non-null, default empty array

	// Authorization
	Role   string `db:"role" json:"role" validate:"oneof=volunteer organizer admin"`
	Status string `db:"status" json:"status" validate:"oneof=pending active suspended deleted"`

	// Verification
	EmailVerified           bool       `db:"email_verified" json:"email_verified"`
	PhoneVerified           bool       `db:"phone_verified" json:"phone_verified"`
	VerificationToken       *string    `db:"verification_token" json:"-"`
	VerificationTokenExpiry *time.Time `db:"verification_token_expiry" json:"-"`

	// Password reset
	PasswordResetToken  *string    `db:"password_reset_token" json:"-"`
	PasswordResetExpiry *time.Time `db:"password_reset_expiry" json:"-"`

	// Login tracking
	LastLoginAt *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	LoginCount  int        `db:"login_count" json:"login_count"`

	// Timestamps
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`

	// Metadata field you're missing
	Metadata json.RawMessage `db:"metadata" json:"metadata"` // non-null, default {}
}
