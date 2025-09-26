package config

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
