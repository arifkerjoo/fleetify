package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// App
	Env     string
	AppPort string
	AppName string

	// DB
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration

	// CORS
	AllowedOrigins string

	// Upload
	MaxUploadSize int64
	UploadPath    string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		// Application
		Env:     getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppName: getEnv("APP_NAME", "Fleetify Backend"),

		// Database
		DBHost: getEnv("DB_HOST", "mysql"),
		DBPort: getEnv("DB_PORT", "3306"),
		DBUser: getEnv("DB_USER", "fleetify"),
		DBPass: getEnv("DB_PASSWORD", ""),
		DBName: getEnv("DB_NAME", "fleetify"),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "supersecret"),
		JWTExpiration: parseDuration(getEnv("JWT_EXPIRATION", "24h")),

		// CORS
		AllowedOrigins: getEnv(
			"ALLOWED_ORIGINS",
			"http://localhost:3000",
		),

		// Upload
		MaxUploadSize: parseInt64(
			getEnv("MAX_UPLOAD_SIZE", "10485760"),
		),
		UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// parseDuration parses duration string
func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Invalid duration format for %s, using default 24h", s)
		return 24 * time.Hour
	}
	return duration
}

// parseInt64 parses string to int64
func parseInt64(s string) int64 {
	var result int64
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int64(r-'0')
		}
	}
	if result == 0 {
		return 10485760 // Default 10MB
	}
	return result
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JWTSecret == "your-super-secret-jwt-key-change-in-production" && c.Env == "production" {
		log.Fatal("❌ JWT_SECRET must be changed in production!")
	}

	if c.DBHost == "" || c.DBUser == "" || c.DBName == "" {
		log.Fatal("❌ Database configuration is incomplete!")
	}

	return nil
}

// IsDevelopment checks if app is in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development" || c.Env == "dev"
}

// IsProduction checks if app is in production mode
func (c *Config) IsProduction() bool {
	return c.Env == "production" || c.Env == "prod"
}
