package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
