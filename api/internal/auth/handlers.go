package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pl3lee/webjson/internal/database"
	"github.com/pl3lee/webjson/internal/utils"
)

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type ApiKeyResponse struct {
	ApiKey string `json:"apiKey"`
}

type ApiKeyMetadata struct {
	Hash       string    `json:"hash"`
	CreatedAt  time.Time `json:"createdAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
}

func (cfg *AuthConfig) HandlerGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url, state, err := cfg.getAuthCodeURL()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot generate google url", err)
		return
	}
	isProd := !strings.Contains(cfg.WebBaseURL, "https")

	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/google/callback",
	})

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (cfg *AuthConfig) HandlerGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// check state to see if they match
	stateCookie, err := r.Cookie("oauthstate")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "state cookie not found", err)
		return
	}
	state := stateCookie.Value
	queryState := r.URL.Query().Get("state")
	if state != queryState {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid state", nil)
		return
	}

	authCode := r.URL.Query().Get("code")
	if authCode == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "auth code not found", nil)
		return
	}

	// exchange code for tokens
	token, err := cfg.exchangeCodeForTokenGoogle(authCode)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "token exchange failed", err)
		return
	}

	// get user info using access token
	userInfo, err := getUserInfoGoogle(token.AccessToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to get user info", err)
		return
	}
	log.Printf("%v logged in\n", userInfo)

	// find user in database
	userDb, err := cfg.Db.GetUserByProviderId(r.Context(), userInfo.Sub)
	if err != nil {
		// user doesn't exist in database
		// save user info in database
		userDb, err = cfg.Db.CreateUser(r.Context(), database.CreateUserParams{
			ProviderID: userInfo.Sub,
			Email:      userInfo.Email,
			Name:       userInfo.Name,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error inserting into database", err)
			return
		}
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error creating session token", err)
		return
	}

	session, err := cfg.createSession(r.Context(), sessionToken, userDb.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error creating session", err)
		return
	}
	isProd := !strings.Contains(cfg.WebBaseURL, "https")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: isProd,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	})

	http.Redirect(w, r, cfg.ClientURL+"/app", http.StatusFound)

}

func (cfg *AuthConfig) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "can't get session token from cookie", err)
		return
	}
	token := sessionCookie.Value
	hashBytes := sha256.Sum256([]byte(token))
	sessionId := hex.EncodeToString(hashBytes[:])

	err = cfg.invalidateSession(r.Context(), sessionId)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "can't invalidate session", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *AuthConfig) HandlerGetMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserIDContextKey).(uuid.UUID)

	user, err := cfg.Db.GetUserById(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func (cfg *AuthConfig) HandlerCreateApiKey(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserIDContextKey).(uuid.UUID)

	apiKey, err := cfg.createApiKey(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot create api key", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, ApiKeyResponse{
		ApiKey: apiKey,
	})
}

func (cfg *AuthConfig) HandlerGetAllApiKeys(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserIDContextKey).(uuid.UUID)

	allApiKeys, err := cfg.Db.GetAllApiKeys(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot get api keys", err)
		return
	}
	var response []ApiKeyMetadata

	for _, apiKey := range allApiKeys {
		response = append(response, ApiKeyMetadata{
			Hash:       apiKey.KeyHash,
			CreatedAt:  apiKey.CreatedAt,
			LastUsedAt: apiKey.LastUsedAt,
		})
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (cfg *AuthConfig) HandlerDeleteApiKey(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserIDContextKey).(uuid.UUID)

	keyHash := chi.URLParam(r, "keyHash")

	apiKey, err := cfg.Db.GetUserFromApiKeyHash(r.Context(), keyHash)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "cannot get api key from db", err)
		return
	}

	if apiKey.UserID != userId {
		utils.RespondWithError(w, http.StatusUnauthorized, "api key does not belong to user", err)
		return
	}

	err = cfg.Db.DeleteApiKey(r.Context(), keyHash)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot delete api key from db", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}
