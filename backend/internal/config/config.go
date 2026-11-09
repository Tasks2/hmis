package config

import (
	"fmt"
	"os"
)

// MinJWTSecretLength is the shortest JWT_SECRET the server accepts.
// HS256 needs at least 256 bits (32 bytes) of key to be secure.
const MinJWTSecretLength = 32

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	if len(jwtSecret) < MinJWTSecretLength {
		return Config{}, fmt.Errorf(
			"JWT_SECRET must be at least %d characters (generate one with: openssl rand -base64 48)",
			MinJWTSecretLength,
		)
	}

	return Config{
		Port:        port,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
	}, nil
}
