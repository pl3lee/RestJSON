package auth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pl3lee/webjson/internal/utils"
)

func (cfg *AuthConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie("jwt")
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Can't get jwt from cookie", err)
			return
		}

		userId, err := ValidateJWT(jwtCookie.Value, cfg.Secret)
		if err != nil {
			// Check if token is expired
			if err == jwt.ErrTokenExpired {
				// Try refresh flow
				if err := cfg.verifyRefreshToken(w, r); err != nil {
					utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
					return
				}
				// After successful refresh, redirect to same endpoint
				http.Redirect(w, r, r.URL.Path, http.StatusTemporaryRedirect)
				return
			}
			// Handle non-expiration errors
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token", err)
			return
		}

		// Valid token, proceed with request
		// Add user ID to context
		ctx := context.WithValue(r.Context(), "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
