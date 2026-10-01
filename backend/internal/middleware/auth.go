package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

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
				http.Error(
					w,
					"missing authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(
					w,
					"invalid authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			token, err := jwt.Parse(
				parts[1],
				func(token *jwt.Token) (interface{}, error) {
					if token.Method != jwt.SigningMethodHS256 {
						return nil, errors.New("unexpected signing method")
					}
					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(
					w,
					"invalid or expired token",
					http.StatusUnauthorized,
				)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(
					w,
					"invalid token claims",
					http.StatusUnauthorized,
				)
				return
			}

			userID, _ := claims["sub"].(string)
			role, _ := claims["role"].(string)

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
