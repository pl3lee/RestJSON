package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pl3lee/webjson/internal/database"
)

func generateSessionToken() (string, error) {
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

func (cfg *AuthConfig) createSession(ctx context.Context, token string, userId uuid.UUID) (database.UserSession, error) {
	hashBytes := sha256.Sum256([]byte(token))
	hash := hex.EncodeToString(hashBytes[:])
	session, err := cfg.Db.StoreUserSession(ctx, database.StoreUserSessionParams{
		ID:        hash,
		UserID:    userId,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // 30 days
	})
	if err != nil {
		return database.UserSession{}, fmt.Errorf("createSession: cannot insert session into db: %w", err)
	}
	return session, nil
}

func (cfg *AuthConfig) validateSessionToken(ctx context.Context, token string) (database.UserSession, database.User, error) {
	hashBytes := sha256.Sum256([]byte(token))
	sessionId := hex.EncodeToString(hashBytes[:])

	session, err := cfg.Db.GetSession(ctx, sessionId)
	if err != nil {
		return database.UserSession{}, database.User{}, fmt.Errorf("validateSessionToken: cannot find session in db: %w", err)
	}
	user, err := cfg.Db.GetUserById(ctx, session.UserID)
	if err != nil {
		return database.UserSession{}, database.User{}, fmt.Errorf("validateSessionToken: cannot get user from session id: %w", err)
	}

	// expired
	if session.ExpiresAt.Before(time.Now()) {
		if err := cfg.Db.InvalidateSession(ctx, sessionId); err != nil {
			return database.UserSession{}, database.User{}, fmt.Errorf("validateSessionToken: cannot invalidate session: %w", err)
		}
		return database.UserSession{}, database.User{}, fmt.Errorf("validateSessionToken: token expired")
	}
	// expires in 15 days, extend session
	if session.ExpiresAt.Before(time.Now().Add(time.Hour * 24 * 15)) {
		session, err = cfg.Db.UpdateSession(ctx, database.UpdateSessionParams{
			ID:        session.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // 30 days
		})
		if err != nil {
			return database.UserSession{}, database.User{}, fmt.Errorf("validateSessionToken: cannot extend session: %w", err)
		}
	}

	return session, user, nil
}

func (cfg *AuthConfig) invalidateSession(ctx context.Context, sessionId string) error {
	if err := cfg.Db.InvalidateSession(ctx, sessionId); err != nil {
		return fmt.Errorf("invalidateSession: cannot invalidate session: %w", err)
	}
	return nil
}

//
// func (cfg *AuthConfig) invalidateAllSessions(ctx context.Context, userId uuid.UUID) error {
// 	if err := cfg.Db.InvalidateAllSessions(ctx, userId); err != nil {
// 		return fmt.Errorf("invalidateAllSessions: cannot invalidate all sessions: %w", err)
// 	}
// 	return nil
// }
