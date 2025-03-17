package auth

import (
	"github.com/pl3lee/restjson/internal/database"
	"github.com/redis/go-redis/v9"
)

type AuthConfig struct {
	Db                 *database.Queries
	BaseURL            string
	GoogleClientID     string
	GoogleClientSecret string
	ClientURL          string
	Rdb                *redis.Client
}
