package middleware

import (
	"context"
	"net/http"
	"strings"
	"testsmth/internal/auth"
)

type contextKey string

const playerIDKey contextKey = "playerID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "authrization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		playerID, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "ivalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), playerIDKey, playerID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PlayerIDFromContext(ctx context.Context) (int64, bool) {
	playerID, ok := ctx.Value(playerIDKey).(int64)
	return playerID, ok
}
