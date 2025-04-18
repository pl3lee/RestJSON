package jsonfile

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/s3util"
	"github.com/pl3lee/restjson/internal/utils"
)

type contextKey string

const FileMetadataContextKey contextKey = "fileId"
const FileContentContextKey contextKey = "fileContent"
const ResourceDataContextKey contextKey = "resourceData"
const ResourceArrayContextKey contextKey = "resourceArray"

// JsonFileMiddleware ensures that the user has access to the requested JSON file.
// This middleware depends on the authMiddleware to run first, which sets the user ID in the context.
func (cfg *JsonConfig) JsonFileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
		fileIdStr := chi.URLParam(r, "fileId")
		fileId, err := uuid.Parse(fileIdStr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "file id not valid", err)
			return
		}

		jsonFileMetadata, err := cfg.Db.GetJsonFile(r.Context(), fileId)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "json file does not exist", err)
			return
		}

		if jsonFileMetadata.UserID != userId {
			utils.RespondWithError(w, http.StatusUnauthorized, "file does not belong to user", nil)
			return
		}

		// valid file id and file belongs to user, proceed with request
		// add file ID to context
		ctx := context.WithValue(r.Context(), FileMetadataContextKey, jsonFileMetadata)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *JsonConfig) JsonFileContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
		fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
		fileContents, err := s3util.GetJsonFromS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error getting json from s3", err)
			return
		}

		ctx := context.WithValue(r.Context(), FileContentContextKey, fileContents)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *JsonConfig) ResourceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileContentsMap := r.Context().Value(FileContentContextKey).(map[string]any)

		resource := chi.URLParam(r, "resource")

		resourceData, ok := fileContentsMap[resource]
		if !ok {
			utils.RespondWithError(w, http.StatusNotFound, "resource not found", nil)
			return
		}

		ctx := context.WithValue(r.Context(), ResourceDataContextKey, resourceData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *JsonConfig) ResourceArrayMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items, ok := r.Context().Value(ResourceDataContextKey).([]any)
		if !ok {
			utils.RespondWithError(w, http.StatusBadRequest, "resource is not an array", nil)
			return
		}

		ctx := context.WithValue(r.Context(), ResourceArrayContextKey, items)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
