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
}

// ServerConfig holds server-related settings (port, environment).
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

// DatabaseConfig holds database connection details.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
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
	v.BindEnv("auth.jwtsecret", "AUTH_JWTSECRET")
	v.BindEnv("auth.issuer", "AUTH_ISSUER")
	v.BindEnv("payment.provider", "PAYMENT_PROVIDER")
	v.BindEnv("payment.apikey", "PAYMENT_APIKEY")

	// 4. Set sensible defaults for development (only used if not overridden).
	v.SetDefault("server.port", 4000)
	v.SetDefault("server.env", "development")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.name", "appdb")
	v.SetDefault("auth.issuer", "myapp")
	v.SetDefault("payment.provider", "stripe")

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
