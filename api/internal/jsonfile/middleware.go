package jsonfile

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/utils"
)

type contextKey string

const FileMetadataContextKey contextKey = "fileId"
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

func (cfg *JsonConfig) ResourceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
		fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
		resource := chi.URLParam(r, "resource")

		resourceData, err := cfg.getResourceFromS3File(r.Context(), userId, fileMetadata.ID, resource)
		if err != nil {
			switch err.(type) {
			case *InternalServerError:
				utils.RespondWithError(w, http.StatusInternalServerError, "internal server error", err)
			case *NotFoundError:
				utils.RespondWithError(w, http.StatusNotFound, "resource not found", err)
			default:
				utils.RespondWithError(w, http.StatusInternalServerError, "unknown error occurred", err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), ResourceDataContextKey, *resourceData)
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
