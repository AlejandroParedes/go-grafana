package config

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
	Sentry   SentryConfig   `json:"sentry"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode"`
}

// LoggingConfig holds logging-specific configuration
type LoggingConfig struct {
	Level string `json:"level"`
}

// SentryConfig holds Sentry-specific configuration
type SentryConfig struct {
	DSN string `json:"dsn"`
}

// NewConfig creates a new configuration instance with environment-based values
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "go_grafana"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		Sentry: SentryConfig{
			DSN: getEnv("SENTRY_DSN", ""),
		},
	}
}

// getEnv retrieves an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv retrieves an environment variable as a duration with a fallback default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		// Try parsing as seconds if it's just a number
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.DBName +
		" sslmode=" + c.Database.SSLMode
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}

// LogConfig logs the configuration (without sensitive data)
func (c *Config) LogConfig(logger *zap.Logger) {
	logger.Info("Configuration loaded",
		zap.String("server_port", c.Server.Port),
		zap.String("db_host", c.Database.Host),
		zap.String("db_port", c.Database.Port),
		zap.String("db_name", c.Database.DBName),
		zap.String("log_level", c.Logging.Level),
		zap.String("sentry_dsn", c.Sentry.DSN),
	)
}
