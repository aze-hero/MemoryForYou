package config

import "os"

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	FrontendURL  string
	OpenAIKey    string
	OpenAIModel  string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("API_PORT", "3001"),
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://timecapsule:timecapsule@localhost:5432/timecapsule?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
		OpenAIModel: getEnv("OPENAI_MODEL", "gpt-4o"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
