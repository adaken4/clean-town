package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/adaken4/clean-town/internal/auth"
	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"github.com/wneessen/go-mail"
	"golang.org/x/crypto/bcrypt"
)

// AuthService provides all authentication and authorization operations
// such as registration, login, logout, token refresh, and email verification.
// It depends on a UserRepository for persistence and a JWT-based token system.
type AuthService struct {
	Config    *config.Config
	UserRepo  models.UserRepository
	Blacklist *auth.InMemoryBlacklist
	Logger    *slog.Logger
}

// HashPassword hashes a plain password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CheckPassword compares a plain password with the stored hash
func CheckPassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

// GenerateEmailVerificationToken creates a signed JWT for verifying user emails.
// The token contains the email address, token type, and expiry time (30 minutes).
func GenerateEmailVerificationToken(signingKey []byte, email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"type":  "email_verification",
		"exp":   time.Now().Add(30 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// SendVerificationEmail sends an HTML email containing a verification link to the user.
// It uses Gmail’s SMTP server for delivery.
func SendVerificationEmail(email, token string, cfg *config.Config) error {
	m := mail.NewMsg()
	if err := m.From(cfg.Email.FromAddress); err != nil {
		return err
	}
	if err := m.To(email); err != nil {
		return err
	}
	m.Subject("Email Verification")

	// Construct the verification URL
	verificationLink := fmt.Sprintf("http://localhost:8080/v1/auth/verify-email?token=%s", token)
	body := fmt.Sprintf(`
		<h2>Verify Your Email</h2>
		<p>Click the link below:</p>
		<a href="%s">Verify Email</a>
	`, verificationLink)

	m.SetBodyString(mail.TypeTextHTML, body)

	// Configure SMTP client
	c, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.Email.FromAddress),
		mail.WithPassword(cfg.Email.AppPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return err
	}

	// Send the email via SMTP
	if err := c.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// Register handles the creation of a new user account
func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	// Clean up user input (trim spaces, normalize casing)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Role = strings.TrimSpace(req.Role)
	req.Town = strings.TrimSpace(req.Town)

	// Hash the password before storing
	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Prepare the user model for persistence
	u := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         req.Role,
		Town:         &req.Town,
	}

	// Generate an email verification token (JWT)
	token, err := GenerateEmailVerificationToken([]byte(s.Config.Auth.JWTSecret), u.Email)
	if err != nil {
		return nil, err
	}

	// Set expiry time for verification token
	expiry := time.Now().Add(31 * time.Minute)

	// Save the new user and verification token
	if err := s.UserRepo.Create(ctx, u, token, expiry); err != nil {
		var pqErr *pq.Error
		// Detect PostgreSQL duplicate key violation (unique email)
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, models.ErrDuplicateEmail
		}
		return nil, err
	}

	// Attempt to send the verification email (log error if sending fails)
	if err := SendVerificationEmail(u.Email, token, s.Config); err != nil {
		s.Logger.Error("failed to send verification email", "error", err)
	}

	return u, nil
}

// Login authenticates a user using their email and password,
// and returns a pair of access + refresh tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	// Retrieve user by email
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", models.ErrInvalidCredentials
		}
	}

	// Generate access token (short-lived)
	access, err := auth.GenerateAccessToken([]byte(s.Config.Auth.JWTSecret), *user)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (longer-lived)
	refresh, err := auth.GenerateRefreshToken([]byte(s.Config.Auth.JWTSecret), *user)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
