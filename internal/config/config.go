package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config is the root configuration struct for the application.
// It groups related settings into logical sub-structs.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Payment  PaymentConfig  `mapstructure:"payment"`
	Email    EmailConfig    `mapstructure:"email"`
}

// ServerConfig holds server-related settings (port, environment).
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

// DatabaseConfig holds database connection details.
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxRetries   int    `mapstructure:"max_retries"`
	RetryBackoff int    `mapstructure:"retry_backoff"` // in secs
	// Connection pool settings
	MaxOpenConns int `mapstructure:"max_open_conns"`
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	MaxIdleTime  int `mapstructure:"max_idle_time"` // in minutes
}

// AuthConfig holds authentication-related settings (JWT secret, issuer).
type AuthConfig struct {
	JWTSecret string `mapstructure:"jwtsecret"`
	Issuer    string `mapstructure:"issuer"`
}

// PaymentConfig holds payment gateway settings (provider, API key).
type PaymentConfig struct {
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"apikey"`
}

type EmailConfig struct {
	FromAddress string `mapstructure:"from"`
	AppPassword string `mapstructure:"app_password"`
}

// LoadConfig loads application configuration by merging settings from:
//  1. An optional .env file (via godotenv).
//  2. An optional config file (YAML/JSON/TOML supported by Viper).
//  3. Environment variables.
//  4. Default values (for development).
//
// Precedence order: ENV VARS > Config File > Defaults
func LoadConfig(path string) (*Config, error) {
	// 1. Load environment variables from .env file if path is provided.
	if path != "" {
		if err := godotenv.Load(path); err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	v := viper.New()

	// 2. Support optional config file (config.yaml/config.json).
	// Looks for "config.{yaml,json,toml}" in the current directory.
	v.SetConfigName("config") // file name without extension
	v.AddConfigPath(".")      // search in project root
	v.SetConfigType("yaml")   // default type = yaml
	if err := v.ReadInConfig(); err != nil {
		// Ignore missing config file, but surface other errors (e.g. parse error).
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// 3. Bind environment variables explicitly.
	// ENV vars take precedence over config file and defaults.
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.env", "SERVER_ENV")
	v.BindEnv("database.host", "DATABASE_HOST")
	v.BindEnv("database.port", "DATABASE_PORT")
	v.BindEnv("database.user", "DATABASE_USER")
	v.BindEnv("database.password", "DATABASE_PASSWORD")
	v.BindEnv("database.name", "DATABASE_NAME")
	v.BindEnv("database.sslmode", "DB_SSL_MODE")
	v.BindEnv("database.max_retries", "MAX_RETRIES")
	v.BindEnv("database.retry_backoff", "RETRY_BACKOFF")
	v.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	v.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	v.BindEnv("database.max_idle_time", "DB_MAX_IDLE_TIME")
	v.BindEnv("auth.jwtsecret", "AUTH_JWTSECRET")
	v.BindEnv("auth.issuer", "AUTH_ISSUER")
	v.BindEnv("payment.provider", "PAYMENT_PROVIDER")
	v.BindEnv("payment.apikey", "PAYMENT_APIKEY")
	v.BindEnv("email.from", "EMAIL_FROM")
	v.BindEnv("email.app_password", "EMAIL_APP_PASSWORD")

	// 4. Set sensible defaults for development (only used if not overridden).
	v.SetDefault("server.port", 4000)
	v.SetDefault("server.env", "development")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.name", "appdb")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_retries", 5)
	v.SetDefault("database.retry_backoff", 1)
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 25)
	v.SetDefault("database.max_idle_time", 15)
	v.SetDefault("auth.issuer", "myapp")
	v.SetDefault("payment.provider", "stripe")
	v.SetDefault("email.from", "no-reply@localhost")

	// 5. Unmarshal merged configuration into strongly typed struct.
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct %w", err)
	}

	// 6. Validate required fields (must be present in at least one source).
	if config.Auth.JWTSecret == "" {
		return nil, fmt.Errorf("missing required auth.jwtsecret")
	}
	if config.Payment.APIKey == "" {
		return nil, fmt.Errorf("missing required payment.apikey")
	}

	return &config, nil
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}
