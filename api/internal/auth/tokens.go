package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "todo-go",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),

		Subject: userID.String(),
	})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := &jwt.RegisteredClaims{}
	// keyFunc is a function that returns the key used to validate the token. Therefore, it returns tokenSecret.

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("Incorrect signing method")
		}
		return []byte(tokenSecret), nil
	}

	// Provide keyFunc to ParseWithClaims so that we can parse the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		claimsStruct,
		keyFunc,
	)
	if err != nil {

		return uuid.UUID{}, err
	}

	// Extract userID from the claims

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	// turns userID string into UUID
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userIDUUID, nil
}

func MakeRefreshToken() (string, error) {
	// Allocate space for 32 bytes (256 bits) of random data
	random := make([]byte, 32)
	// writes random data into the slice
	_, err := rand.Read(random)
	if err != nil {
		return "", fmt.Errorf("error in generating random string")
	}
	// encodes as hex
	randomString := hex.EncodeToString(random)
	return randomString, nil
}

func (cfg *AuthConfig) verifyRefreshToken(w http.ResponseWriter, r *http.Request) error {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return fmt.Errorf("no refresh token")
	}

	// Verify refresh token in database
	token, err := cfg.Db.GetRefreshToken(r.Context(), refreshCookie.Value)
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}

	// Issue new JWT
	newJWT, err := MakeJWT(token.UserID, cfg.Secret, time.Hour)
	if err != nil {
		return fmt.Errorf("error creating new JWT")
	}

	// Set new JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    newJWT,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour),
	})

	return nil
}
