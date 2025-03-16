package jsonfile

import (
	"encoding/json"
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

	fileContents, err := cfg.getFileFromS3(r.Context(), userId, fileMetadata.ID)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot get json file from s3", err)
		return
	}
	var fileContentsMap map[string]any
	if err := json.Unmarshal(fileContents, &fileContentsMap); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot unmarshal json file", err)
		return
	}

	resourceData, ok := fileContentsMap[resource]
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "resource not found", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, resourceData)
}
