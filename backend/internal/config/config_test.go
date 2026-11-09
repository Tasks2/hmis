package config

import (
	"strings"
	"testing"
)

const validSecret = "0123456789abcdef0123456789abcdef"

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		databaseURL string
		jwtSecret   string
		wantErr     string
	}{
		{
			name:        "valid config",
			databaseURL: "postgres://localhost/hmis",
			jwtSecret:   validSecret,
		},
		{
			name:      "missing database URL",
			jwtSecret: validSecret,
			wantErr:   "DATABASE_URL is required",
		},
		{
			name:        "missing JWT secret",
			databaseURL: "postgres://localhost/hmis",
			wantErr:     "JWT_SECRET is required",
		},
		{
			name:        "JWT secret too short",
			databaseURL: "postgres://localhost/hmis",
			jwtSecret:   validSecret[:MinJWTSecretLength-1],
			wantErr:     "at least 32 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", "")
			t.Setenv("DATABASE_URL", tt.databaseURL)
			t.Setenv("JWT_SECRET", tt.jwtSecret)

			cfg, err := Load()

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Port != "8080" {
				t.Errorf("expected default port 8080, got %q", cfg.Port)
			}

			if cfg.JWTSecret != tt.jwtSecret {
				t.Errorf("JWT secret not loaded")
			}
		})
	}
}
