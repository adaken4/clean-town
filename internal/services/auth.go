package services

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/adaken4/clean-town/internal/auth"
	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/models"
	"github.com/golang-jwt/jwt/v5"
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
