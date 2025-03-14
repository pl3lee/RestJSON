package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/pl3lee/webjson/internal/utils"
)

type contextKey string

const UserIDContextKey contextKey = "userId"

func (cfg *AuthConfig) SessionMiddleware(next http.Handler) http.Handler {
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
		ctx := context.WithValue(r.Context(), UserIDContextKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *AuthConfig) ApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get api key from header
		authHeader := r.Header.Get("Authorization")
		words := strings.Split(authHeader, " ")
		if len(words) != 2 || words[0] != "Bearer" {
			utils.RespondWithError(w, http.StatusUnauthorized, "malformed authorization header", nil)
			return
		}
		apiKey := words[1]

		// hash api key
		apiKeyHashBytes := sha256.Sum256([]byte(apiKey))
		apiKeyHash := hex.EncodeToString(apiKeyHashBytes[:])

		// lookup database
		apiKeyEntry, err := cfg.Db.GetUserFromApiKeyHash(r.Context(), apiKeyHash)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "invalid api key", err)
			return
		}
		userId := apiKeyEntry.UserID
		err = cfg.Db.UpdateApiKeyLastUsed(r.Context(), apiKeyHash)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error updating api key last used", err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
