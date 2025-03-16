package jsonfile

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/utils"
)

func (cfg *JsonConfig) HandlerGetResource(w http.ResponseWriter, r *http.Request) {
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
			utils.RespondWithError(w, http.StatusInternalServerError, "unknown error occured", err)
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, *resourceData)
}

func (cfg *JsonConfig) HandlerGetResourceById(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
	resource := chi.URLParam(r, "resource")
	resourceId := chi.URLParam(r, "id")

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

	items, ok := (*resourceData).([]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "resource is not an array, and thus cannot be indexed with id", nil)
		return
	}

	for _, item := range items {
		itemMap, ok := item.(map[string]any)
		if ok {
			id := fmt.Sprintf("%v", itemMap["id"])
			if id == resourceId {
				utils.RespondWithJSON(w, http.StatusOK, item)
				return
			}
		}
	}
	utils.RespondWithError(w, http.StatusNotFound, "resource with particular id not found", nil)
}
