package auth

import "github.com/pl3lee/webjson/internal/database"

type AuthConfig struct {
	Db                 *database.Queries
	BaseURL            string
	GoogleClientID     string
	GoogleClientSecret string
	ClientURL          string
}
