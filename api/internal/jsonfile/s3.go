package jsonfile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func (cfg *JsonConfig) uploadJsonToS3(ctx context.Context, userId uuid.UUID, fileId uuid.UUID, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("uploadJsonToS3: error marshalling data: %w", err)
	}

	tempFile, err := os.CreateTemp("", fileId.String())
	if err != nil {
		return fmt.Errorf("uploadJsonToS3: error creating temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// write json data to temp file
	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("uploadJsonToS3: error writing data to temp file: %w", err)
	}

	// reset file pointer to beginning
	if _, err := tempFile.Seek(0, 0); err != nil {
		return fmt.Errorf("uploadJsonToS3: error resetting file pointer: %w", err)
	}

	if err := cfg.uploadFileToS3(ctx, userId, fileId, tempFile); err != nil {
		return fmt.Errorf("uploadJsonToS3: error uploading file to s3: %w", err)
	}
	return nil
}

func (cfg *JsonConfig) uploadFileToS3(ctx context.Context, userId, fileId uuid.UUID, file io.Reader) error {
	_, err := cfg.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3Bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s.json", userId.String(), fileId.String())),
		Body:        file,
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("uploadFileToS3: error uploading to S3: %w", err)
	}
	return nil
}

func (cfg *JsonConfig) getJsonFromS3(ctx context.Context, userId, fileId uuid.UUID) (any, error) {
	data, err := cfg.getFileFromS3(ctx, userId, fileId)
	if err != nil {
		return nil, fmt.Errorf("getJsonFromS3: error getting file from S3: %w", err)
	}

	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("getJsonFromS3: error unmarshalling json: %w", err)
	}

	return result, nil
}

func (cfg *JsonConfig) getFileFromS3(ctx context.Context, userId, fileId uuid.UUID) ([]byte, error) {
	s3Params := &s3.GetObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s.json", userId.String(), fileId.String())),
	}

	result, err := cfg.S3Client.GetObject(ctx, s3Params)
	if err != nil {
		return nil, fmt.Errorf("getFileFromS3: error getting object from S3: %w", err)
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("getFileFromS3: error reading ojbect body: %w", err)
	}

	return body, nil
}

func (cfg *JsonConfig) deleteFileFromS3(ctx context.Context, userId, fileId uuid.UUID) error {
	s3Params := &s3.DeleteObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s.json", userId.String(), fileId.String())),
	}
	_, err := cfg.S3Client.DeleteObject(ctx, s3Params)
	if err != nil {
		return fmt.Errorf("deleteFileFromS3: error deleting object frmo S3: %w", err)
	}
	return nil
}
