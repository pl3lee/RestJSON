package config

import "github.com/pl3lee/webjson/internal/database"

type Config struct {
	Port               string
	ClientURL          string
	DbUrl              string
	Secret             string
	WebBaseURL         string
	GoogleClientID     string
	GoogleClientSecret string
	Db                 *database.Queries
}
