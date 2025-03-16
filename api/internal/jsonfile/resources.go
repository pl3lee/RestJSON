package jsonfile

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (cfg *JsonConfig) getResourceFromS3File(ctx context.Context, userId, fileId uuid.UUID, resource string) (*any, error) {
	fileContents, err := cfg.getFileFromS3(ctx, userId, fileId)

	if err != nil {
		return nil, &InternalServerError{
			Message: fmt.Sprintf("getResourceFromS3File: cannot get file contents: %v", err),
		}
	}
	var fileContentsMap map[string]any
	if err := json.Unmarshal(fileContents, &fileContentsMap); err != nil {
		return nil, &InternalServerError{
			Message: fmt.Sprintf("getResourceFromS3File: cannot unmarshal json file: %v", err),
		}
	}

	resourceData, ok := fileContentsMap[resource]
	if !ok {
		return nil, &NotFoundError{
			Message: "getResourceFromS3File: resource not found",
		}
	}
	return &resourceData, nil
}
