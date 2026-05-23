package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	JWTAccessExpiry    int
	JWTRefreshExpiry   int
	FrontendURL        string
	OpenAIKey          string
	OpenAIModel        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("API_PORT", "3001"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgresql://timecapsule:timecapsule@localhost:5432/timecapsule?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "change-me"),
		JWTAccessExpiry:    getEnvInt("JWT_ACCESS_EXPIRY", 3600),
		JWTRefreshExpiry:   getEnvInt("JWT_REFRESH_EXPIRY", 2592000),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		OpenAIKey:          getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:        getEnv("OPENAI_MODEL", "gpt-4o"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL:  getEnv("GITHUB_REDIRECT_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var result int
		if _, err := fmt.Sscanf(v, "%d", &result); err == nil {
			return result
		}
	}
	return fallback
}
