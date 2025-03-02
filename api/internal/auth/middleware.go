package auth

import (
	"context"
	"net/http"

	"github.com/pl3lee/webjson/internal/utils"
)

type contextKey string

const userIDKey contextKey = "user_id"

func (cfg *AuthConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "can't get session token from cookie", err)
			return
		}

		_, user, err := cfg.validateSessionToken(r.Context(), sessionCookie.Value)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "cannot validate session", err)
			return
		}

		// valid token, proceed with request
		// add user ID to context
		ctx := context.WithValue(r.Context(), userIDKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
