package auth

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	S3Bucket           string
	S3Region           string
	S3Client           *s3.Client
}
