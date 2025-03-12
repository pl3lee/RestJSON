package jsonfile

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pl3lee/webjson/internal/auth"
	"github.com/pl3lee/webjson/internal/utils"
)

type contextKey string

const FileIDContextKey contextKey = "file_id"

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
		ctx := context.WithValue(r.Context(), FileIDContextKey, fileId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
