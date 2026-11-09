package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tasks2/hmis/internal/response"
	"github.com/Tasks2/hmis/internal/services"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func signToken(t *testing.T, method jwt.SigningMethod, secret string, claims jwt.MapClaims) string {
	t.Helper()

	token, err := jwt.NewWithClaims(method, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return token
}

func validClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"sub":  "user-123",
		"role": "PATIENT",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
}

// Regression tests: each of these tokens must be rejected with exactly one
// 401 JSON response, and the protected handler must never run.
func TestJWTAuth_RejectsInvalidTokens(t *testing.T) {
	expired := validClaims()
	expired["exp"] = time.Now().Add(-time.Minute).Unix()

	noRole := validClaims()
	delete(noRole, "role")

	tests := []struct {
		name     string
		token    string
		wantCode string
	}{
		{
			name:     "HS384 token signed with the right secret",
			token:    signToken(t, jwt.SigningMethodHS384, testSecret, validClaims()),
			wantCode: "INVALID_TOKEN",
		},
		{
			name:     "HS512 token signed with the right secret",
			token:    signToken(t, jwt.SigningMethodHS512, testSecret, validClaims()),
			wantCode: "INVALID_TOKEN",
		},
		{
			name:     "token signed with the wrong secret",
			token:    signToken(t, jwt.SigningMethodHS256, "another-secret", validClaims()),
			wantCode: "INVALID_TOKEN",
		},
		{
			name:     "expired token",
			token:    signToken(t, jwt.SigningMethodHS256, testSecret, expired),
			wantCode: "INVALID_TOKEN",
		},
		{
			name:     "token without a role",
			token:    signToken(t, jwt.SigningMethodHS256, testSecret, noRole),
			wantCode: "INVALID_TOKEN_CLAIMS",
		},
		{
			name:     "not a JWT",
			token:    "not-a-token",
			wantCode: "INVALID_TOKEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := JWTAuth(testSecret)(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("handler should not be called")
				}),
			)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d", rec.Code)
			}

			// Decoding the whole body as one object fails if two error
			// responses were written, which the old code did.
			var body response.ErrorResponse
			dec := json.NewDecoder(rec.Body)
			if err := dec.Decode(&body); err != nil {
				t.Fatalf("body is not a JSON error: %v", err)
			}
			if dec.More() {
				t.Fatal("more than one response was written")
			}

			if body.Error.Code != tt.wantCode {
				t.Errorf("expected code %s, got %s", tt.wantCode, body.Error.Code)
			}
		})
	}
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	handler := JWTAuth(
		"test-secret",
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("handler should not be called")
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	jwtService := services.NewJWTService("test-secret")

	token, err := jwtService.GenerateToken(
		"user-123",
		"PATIENT",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	handler := JWTAuth(
		"test-secret",
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(string)
			if !ok {
				t.Fatal("user ID missing from context")
			}

			if userID != "user-123" {
				t.Fatalf("expected user-123, got %s", userID)
			}

			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
