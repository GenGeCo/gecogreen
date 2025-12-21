package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	Port    string
	AppEnv  string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// JWT
	JWTSecret string

	// Stripe
	StripeSecretKey string
	StripeWebhookSecret string

	// Cloudflare R2
	R2AccountID   string
	R2AccessKeyID string
	R2SecretKey   string
	R2BucketName  string
	R2PublicURL   string

	// SMTP (Email)
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	FrontendURL  string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file in development
	_ = godotenv.Load()

	return &Config{
		// Server
		Port:    getEnv("PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgres://gecogreen:gecogreen_dev_2024@localhost:5432/gecogreen?sslmode=disable"),

		// Redis
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "dev_jwt_secret_change_in_production"),

		// Stripe
		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),

		// Cloudflare R2
		R2AccountID:   getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID: getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretKey:   getEnv("R2_SECRET_KEY", ""),
		R2BucketName:  getEnv("R2_BUCKET_NAME", "gecogreen-uploads"),
		R2PublicURL:   getEnv("R2_PUBLIC_URL", ""),

		// SMTP
		SMTPHost:     getEnv("SMTP_HOST", "smtps.aruba.it"),
		SMTPPort:     getEnv("SMTP_PORT", "465"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@gecogreen.com"),
		FrontendURL:  getEnv("FRONTEND_URL", "https://gecogreen.com"),
	}
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
