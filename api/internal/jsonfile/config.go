package jsonfile

import "github.com/pl3lee/webjson/internal/database"

type JsonConfig struct {
	Db         *database.Queries
	WebBaseURL string
	ClientURL  string
}
