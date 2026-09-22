package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

type Config struct {
	Port               string
	DatabaseURL        string
	GoogleClientID     string
	CORSAllowedOrigins []string
	SessionTTL         time.Duration
	AppEnv             string
}

func (c Config) IsProduction() bool { return c.AppEnv == EnvProduction }

func (c Config) CookieSecure() bool { return c.IsProduction() }

func Load() (Config, error) {
	ttlHours, err := lookupInt("SESSION_TTL_HOURS", 720)
	if err != nil {
		return Config{}, err
	}

	appEnv := lookupString("APP_ENV", EnvDevelopment)
	if appEnv != EnvDevelopment && appEnv != EnvProduction {
		return Config{}, fmt.Errorf("APP_ENV must be %q or %q, got %q", EnvDevelopment, EnvProduction, appEnv)
	}

	cfg := Config{
		Port:               lookupString("PORT", "8080"),
		DatabaseURL:        lookupString("DATABASE_URL", ""),
		GoogleClientID:     lookupString("GOOGLE_CLIENT_ID", ""),
		CORSAllowedOrigins: splitOrigins(lookupString("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),
		SessionTTL:         time.Duration(ttlHours) * time.Hour,
		AppEnv:             appEnv,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.IsProduction() && cfg.GoogleClientID == "" {
		return Config{}, fmt.Errorf("GOOGLE_CLIENT_ID is required when APP_ENV=%s", EnvProduction)
	}
	if len(cfg.CORSAllowedOrigins) == 0 {
		return Config{}, fmt.Errorf("CORS_ALLOWED_ORIGINS must contain at least one origin")
	}
	for _, origin := range cfg.CORSAllowedOrigins {
		if origin == "*" {
			return Config{}, fmt.Errorf("CORS_ALLOWED_ORIGINS must not contain %q because credentials are enabled", "*")
		}
	}

	return cfg, nil
}

func lookupString(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return fallback
}

func lookupInt(key string, fallback int) (int, error) {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	if parsed <= 0 {
		return 0, fmt.Errorf("%s must be positive, got %d", key, parsed)
	}
	return parsed, nil
}

func splitOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}
