package config

import "os"

type Config struct {
	DatabaseURL     string
	JWTSecret       string
	Port            string
	OpenWebUIURL    string
	OpenWebUIAPIKey string
}

func Load() *Config {
	return &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://complisec:complisec123@postgres:5432/complisec?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		Port:            getEnv("PORT", "8080"),
		OpenWebUIURL:    getEnv("OPENWEBUI_URL", ""),
		OpenWebUIAPIKey: getEnv("OPENWEBUI_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
