package jsonfile

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/redis/go-redis/v9"
)

type JsonConfig struct {
	Db            *database.Queries
	BaseURL       string
	ClientURL     string
	S3Bucket      string
	S3Region      string
	S3Client      *s3.Client
	Rdb           *redis.Client
	FreeFileLimit int
	ProFileLimit  int
}
