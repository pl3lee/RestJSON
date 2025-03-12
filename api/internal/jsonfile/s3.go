package jsonfile

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

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
