package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pl3lee/webjson/internal/database"
	"github.com/pl3lee/webjson/internal/utils"
)

func (cfg *AuthConfig) HandlerGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// TODO: generate google url

	// TODO: redirect user to google url
}

func (cfg *AuthConfig) HandlerGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// TODO: Change this to get the code from query params
	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Exchange code for tokens
	token, err := cfg.exchangeCodeForTokenGoogle(req.Code)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Token exchange failed", err)
		return
	}

	// Get user info using access token
	userInfo, err := getUserInfoGoogle(token.AccessToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get user info", err)
		return
	}
	fmt.Printf("%v logged in", userInfo)

	// find user in database
	userDb, err := cfg.Db.GetUserByProviderId(r.Context(), userInfo.Sub)
	if err != nil {
		// user doesn't exist in database
		// Save user info in database
		userDb, err = cfg.Db.CreateUser(r.Context(), database.CreateUserParams{
			ProviderID: userInfo.Sub,
			Email:      userInfo.Email,
			Name:       userInfo.Name,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error inserting into database", err)
			return
		}
	}

	var expiration time.Duration = time.Hour
	// Create JWT
	jwtToken, err := MakeJWT(userDb.ID, cfg.Secret, expiration)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
		return
	}

	refreshToken, err := MakeRefreshToken()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}
	_, err = cfg.Db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userDb.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error storing refresh token", err)
		return
	}

	// Set jwt cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Path:     "/",
		Domain:   "", // Empty for same-origin
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(expiration),
	})

	// Set refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24 * 60),
	})

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"user":   userDb,
	})

}

func (cfg *AuthConfig) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "logged out",
	})
}

func (cfg *AuthConfig) HandlerGetMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(uuid.UUID)

	user, err := cfg.Db.GetUserById(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}
