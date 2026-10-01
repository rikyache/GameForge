package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"testsmth/internal/auth"
)

type contextKey string

const (
	playerIDKey       contextKey = "playerID"
	tokenIDKey        contextKey = "tokenID"
	tokenExpiresAtKey contextKey = "tokenExpiresAt"
)

func Auth(blacklist auth.TokenBlacklist) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			tokenData, err := auth.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			revoked, err := blacklist.IsRevoked(r.Context(), tokenData.TokenID)
			if err != nil {
				http.Error(w, "failed to check token", http.StatusInternalServerError)
				return
			}

			if revoked {
				http.Error(w, "token is revoked", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, playerIDKey, tokenData.PlayerID)
			ctx = context.WithValue(ctx, tokenIDKey, tokenData.TokenID)
			ctx = context.WithValue(ctx, tokenExpiresAtKey, tokenData.ExpiresAt)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func PlayerIDFromContext(ctx context.Context) (int64, bool) {
	playerID, ok := ctx.Value(playerIDKey).(int64)
	return playerID, ok
}

func TokenIDFromContext(ctx context.Context) (string, bool) {
	tokenID, ok := ctx.Value(tokenIDKey).(string)
	return tokenID, ok
}

func TokenExpiresAtFromContext(ctx context.Context) (time.Time, bool) {
	expiresAt, ok := ctx.Value(tokenExpiresAtKey).(time.Time)
	return expiresAt, ok
}
