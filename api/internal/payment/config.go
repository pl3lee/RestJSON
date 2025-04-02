package payment

import "github.com/pl3lee/restjson/internal/database"

type JsonConfig struct {
	Db        *database.Queries
	BaseURL   string
	ClientURL string
}
