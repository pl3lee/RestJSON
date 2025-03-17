package jsonfile

import (
	"encoding/json"
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

func (cfg *JsonConfig) HandlerGetResourceItem(w http.ResponseWriter, r *http.Request) {
	items, ok := r.Context().Value(ResourceDataContextKey).([]any)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve items from context", nil)
		return
	}
	resourceId := chi.URLParam(r, "id")

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

func (cfg *JsonConfig) HandlerCreateResourceItem(w http.ResponseWriter, r *http.Request) {
	// TODO: get whole json file and modify it
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
	items, ok := r.Context().Value(ResourceArrayContextKey).([]any)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve items from context", nil)
		return
	}
	var newResource map[string]any
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&newResource); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	items = append(items, newResource)

}
