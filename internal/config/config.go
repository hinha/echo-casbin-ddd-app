package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	APIKey   APIKeyConfig
}

// ServerConfig holds all server related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds all database related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds all JWT related configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// APIKeyConfig holds all API Key related configuration
type APIKeyConfig struct {
	HeaderName string
}

// loadEnvFiles loads environment variables from .env* files
func loadEnvFiles() error {
	// Find all .env* files in the current directory
	matches, err := filepath.Glob(".env*")
	if err != nil {
		return err
	}

	// Load each env file
	for _, match := range matches {
		if err := godotenv.Load(match); err != nil {
			return err
		}
	}

	return nil
}

// NewConfig creates a new Config
func NewConfig() *Config {
	// Load environment variables from env.* files
	if err := loadEnvFiles(); err != nil {
		// Log the error but continue with default values
		// This allows the application to start even if env files are not found
		// or cannot be loaded
		println("Warning: Failed to load environment files:", err.Error())
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "pass"),
			DBName:   getEnv("DB_NAME", "auth"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		APIKey: APIKeyConfig{
			HeaderName: getEnv("API_KEY_HEADER", "X-API-Key"),
		},
	}
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an environment variable as an integer or a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Helper function to get an environment variable as a duration or a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
