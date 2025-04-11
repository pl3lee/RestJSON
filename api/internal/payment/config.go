package payment

import (
	"github.com/pl3lee/restjson/internal/database"
	"github.com/redis/go-redis/v9"
)

type PaymentConfig struct {
	Db        *database.Queries
	BaseURL   string
	ClientURL string
	Rdb       *redis.Client
}
