package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tasks2/hmis/internal/services"
)

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
