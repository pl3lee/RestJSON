package s3util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func getFileFromS3(ctx context.Context, client *s3.Client, bucket string, userId uuid.UUID, fileId uuid.UUID) ([]byte, error) {
	fmt.Println("Getting data from S3")
	s3Params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s.json", userId.String(), fileId.String())),
	}

	result, err := client.GetObject(ctx, s3Params)
	if err != nil {
		return nil, fmt.Errorf("getFileFromS3: error getting object from S3: %w", err)
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("getFileFromS3: error reading object body: %w", err)
	}

	return body, nil
}

func uploadFileToS3(ctx context.Context, client *s3.Client, bucket string, userId uuid.UUID, fileId uuid.UUID, file io.Reader) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s.json", userId.String(), fileId.String())),
		Body:        file,
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("uploadFileToS3: error uploading to S3: %w", err)
	}
	return nil
}

func GetJsonFromS3(ctx context.Context, client *s3.Client, rdb *redis.Client, bucket string, userId uuid.UUID, fileId uuid.UUID) (any, error) {
	// get from cache first
	cacheKey := fmt.Sprintf("json:%s:%s", userId.String(), fileId.String())
	cachedData, err := rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		// cache hit
		var result any
		if err := json.Unmarshal(cachedData, &result); err != nil {
			fmt.Printf("getJsonFromS3: failed to unmarshal cached json: %v\n", err)
		} else {
			return result, nil
		}
	} else if err != redis.Nil {
		// redis error, not just cache miss
		fmt.Printf("GetJsonFromS3: redis error: %v\n", err)
	}

	// cache miss or error
	data, err := getFileFromS3(ctx, client, bucket, userId, fileId)
	if err != nil {
		return nil, fmt.Errorf("GetJsonFromS3: error getting file from S3: %w", err)
	}

	// cache it
	rdb.Set(ctx, cacheKey, data, 24*time.Hour)

	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("GetJsonFromS3: error unmarshalling json: %w", err)
	}

	return result, nil
}

func UploadJsonToS3(ctx context.Context, client *s3.Client, rdb *redis.Client, bucket string, userId uuid.UUID, fileId uuid.UUID, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("UploadJsonToS3: error marshalling data: %w", err)
	}

	tempFile, err := os.CreateTemp("", fileId.String())
	if err != nil {
		return fmt.Errorf("UploadJsonToS3: error creating temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// write json data to temp file
	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("UploadJsonToS3: error writing data to temp file: %w", err)
	}

	// reset file pointer to beginning
	if _, err := tempFile.Seek(0, 0); err != nil {
		return fmt.Errorf("UploadJsonToS3: error resetting file pointer: %w", err)
	}

	if err := uploadFileToS3(ctx, client, bucket, userId, fileId, tempFile); err != nil {
		return fmt.Errorf("UploadJsonToS3: error uploading file to s3: %w", err)
	}

	// cache json file
	cacheKey := fmt.Sprintf("json:%s:%s", userId.String(), fileId.String())
	err = rdb.Set(ctx, cacheKey, data, 24*time.Hour).Err()
	if err != nil {
		fmt.Printf("UploadJsonToS3: failed to cache JSON: %v\n", err)
	}
	return nil
}

func DeleteFileFromS3(ctx context.Context, client *s3.Client, rdb *redis.Client, bucket string, userId uuid.UUID, fileId uuid.UUID) error {
	s3Params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s.json", userId.String(), fileId.String())),
	}
	_, err := client.DeleteObject(ctx, s3Params)
	if err != nil {
		return fmt.Errorf("deleteFileFromS3: error deleting object from S3: %w", err)
	}

	// delete from cache
	cacheKey := fmt.Sprintf("json:%s:%s", userId.String(), fileId.String())
	rdb.Del(ctx, cacheKey)

	return nil
}
