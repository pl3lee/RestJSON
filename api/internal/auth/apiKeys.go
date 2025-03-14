package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/pl3lee/webjson/internal/database"
)

func (cfg *AuthConfig) createApiKey(ctx context.Context, userId uuid.UUID) (string, error) {
	// Allocate space for 32 bytes (256 bits) of random data
	random := make([]byte, 32)
	_, err := rand.Read(random)
	if err != nil {
		return "", fmt.Errorf("error in generating random string")
	}
	// this will be given to the user
	apiKey := hex.EncodeToString(random)

	hashBytes := sha256.Sum256([]byte(apiKey))
	apiKeyHash := hex.EncodeToString(hashBytes[:])

	_, err = cfg.Db.CreateApiKey(ctx, database.CreateApiKeyParams{
		UserID:  userId,
		KeyHash: apiKeyHash,
	})
	if err != nil {
		return "", fmt.Errorf("error in storing api key to database")
	}
	return apiKey, nil

}
