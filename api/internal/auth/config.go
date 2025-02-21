package auth

import "github.com/pl3lee/webjson/internal/database"

type AuthConfig struct {
	Db                 *database.Queries
	AuthSecret         string
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
}
