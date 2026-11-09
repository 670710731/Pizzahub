package mw

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/670710731/Pizzahub/backend/auth"
)

type contextKey string

const userContextKey contextKey = "authenticated-user"

type Identity struct {
	UserID int64
	Role   string
}

func CurrentUser(ctx context.Context) (Identity, bool) {
	identity, ok := ctx.Value(userContextKey).(Identity)
	return identity, ok
}

func RequireAuth(tokens *auth.TokenService, adminOnly bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		scheme, token, ok := strings.Cut(header, " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
			writeAuthError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		claims, err := tokens.Verify(token)
		if err != nil {
			writeAuthError(w, http.StatusUnauthorized, "invalid or expired access token")
			return
		}
		if adminOnly && claims.Role != "admin" {
			writeAuthError(w, http.StatusForbidden, "admin access required")
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, Identity{UserID: claims.UserID, Role: claims.Role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
