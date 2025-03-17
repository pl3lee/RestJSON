package jsonfile

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/utils"
)

type CreateJsonRequest struct {
	FileName string `json:"fileName"`
}

type RenameJsonRequest struct {
	FileName string `json:"fileName"`
}

type JsonMetadataResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	FileName   string    `json:"fileName"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

func (cfg *JsonConfig) HandlerCreateJson(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)

	var createReq CreateJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	defer r.Body.Close()

	if createReq.FileName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "file name cannot be empty", nil)
		return
	}

	// create json file
	fileId := uuid.New()

	emptyJson := map[string]any{}
	err := cfg.uploadJsonToS3(r.Context(), userId, fileId, emptyJson)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error uploading empty JSON to s3", err)
		return
	}
	file, err := cfg.Db.CreateNewJson(r.Context(), database.CreateNewJsonParams{
		ID:       fileId,
		UserID:   userId,
		FileName: createReq.FileName,
		Url:      fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/%s.json", cfg.S3Bucket, cfg.S3Region, userId.String(), fileId.String()),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "cannot create new json", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, file)
}

func (cfg *JsonConfig) HandlerUpdateJson(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)

	var jsonData any
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid JSON data", err)
		return
	}
	defer r.Body.Close()

	if err := cfg.uploadJsonToS3(r.Context(), userId, fileMetadata.ID, jsonData); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error uploading JSON to s3", err)
		return
	}

	fileContents, err := cfg.getJsonFromS3(r.Context(), userId, fileMetadata.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot get updated json file from s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, fileContents)
}

func (cfg *JsonConfig) HandlerGetJsonMetadata(w http.ResponseWriter, r *http.Request) {
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)

	response := JsonMetadataResponse{
		ID:       fileMetadata.ID,
		UserID:   fileMetadata.UserID,
		FileName: fileMetadata.FileName,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)

}

func (cfg *JsonConfig) HandlerGetJson(w http.ResponseWriter, r *http.Request) {
	fileContents := r.Context().Value(FileContentContextKey)
	utils.RespondWithJSON(w, http.StatusOK, fileContents)
}

func (cfg *JsonConfig) HandlerGetJsonFiles(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)

	jsonFiles, err := cfg.Db.GetJsonFiles(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error getting json files", err)
		return
	}

	var jsonFilesResponse []JsonMetadataResponse
	for _, file := range jsonFiles {
		fileMetadata := JsonMetadataResponse{
			ID:         file.ID,
			UserID:     file.UserID,
			FileName:   file.FileName,
			ModifiedAt: file.UpdatedAt,
		}
		jsonFilesResponse = append(jsonFilesResponse, fileMetadata)
	}

	utils.RespondWithJSON(w, http.StatusOK, jsonFilesResponse)
}

func (cfg *JsonConfig) HandlerRenameJsonFile(w http.ResponseWriter, r *http.Request) {
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)

	var renameReq RenameJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&renameReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	defer r.Body.Close()

	if renameReq.FileName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "file name cannot be empty", nil)
		return
	}

	renamedJsonFile, err := cfg.Db.RenameJsonFile(r.Context(), database.RenameJsonFileParams{
		ID:       fileMetadata.ID,
		FileName: renameReq.FileName,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error renaming json file", err)
		return
	}

	response := JsonMetadataResponse{
		ID:       renamedJsonFile.ID,
		UserID:   renamedJsonFile.UserID,
		FileName: renamedJsonFile.FileName,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (cfg *JsonConfig) HandlerDeleteJsonFile(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)

	err := cfg.Db.DeleteJsonFile(r.Context(), fileMetadata.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error deleting json file from db", err)
		return
	}
	if err := cfg.deleteFileFromS3(r.Context(), userId, fileMetadata.ID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error deleting json file from s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}
