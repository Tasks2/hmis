package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Tasks2/hmis/internal/response"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				response.JSONError(
					w,
					http.StatusUnauthorized,
					"UNAUTHORIZED",
					"authorization header is required",
				)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				response.JSONError(
					w,
					http.StatusUnauthorized,
					"INVALID_AUTH_HEADER",
					"authorization header must use Bearer token format",
				)
				return
			}

			token, err := jwt.Parse(
				parts[1],
				func(token *jwt.Token) (interface{}, error) {
					if token.Method != jwt.SigningMethodHS256 {
						response.JSONError(
							w,
							http.StatusUnauthorized,
							"INVALID_TOKEN",
							"invalid token signing method",
						)
					}
					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				response.JSONError(
					w,
					http.StatusUnauthorized,
					"INVALID_TOKEN",
					"token is invalid or expired",
				)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.JSONError(
					w,
					http.StatusUnauthorized,
					"INVALID_TOKEN_CLAIMS",
					"token claims are invalid",
				)
				return
			}

			userID, _ := claims["sub"].(string)
			role, _ := claims["role"].(string)

			if userID == "" || role == "" {
				response.JSONError(
					w,
					http.StatusUnauthorized,
					"INVALID_TOKEN_CLAIMS",
					"token claims are invalid",
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				userID,
			)

			ctx = context.WithValue(
				ctx,
				RoleKey,
				role,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleKey).(string)

			if !ok || role == "" {
				response.JSONError(
					w,
					http.StatusForbidden,
					"FORBIDDEN",
					"role is required",
				)
				return
			}

			if role != requiredRole {
				response.JSONError(
					w,
					http.StatusForbidden,
					"FORBIDDEN",
					"insufficient permissions",
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
