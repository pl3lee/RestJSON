package jsonfile

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/pl3lee/webjson/internal/auth"
	"github.com/pl3lee/webjson/internal/database"
	"github.com/pl3lee/webjson/internal/utils"
)

type CreateJsonRequest struct {
	FileName string `json:"fileName"`
}

type RenameJsonRequest struct {
	FileName string `json:"fileName"`
}

type JsonMetadataResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"userId"`
	FileName string    `json:"fileName"`
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

	// temp file with empty JSON content
	tempFile, err := os.CreateTemp("", fileId.String())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot create temp file", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// write empty JSON content
	if _, err := tempFile.Write([]byte("{}")); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot write to temp file", err)
		return
	}

	// Reset the file pointer to the beginning
	if _, err := tempFile.Seek(0, 0); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot seek temp file", err)
		return
	}

	// upload file to s3
	if err := cfg.uploadFileToS3(r.Context(), userId, fileId, tempFile); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error uploading file to s3", err)
		return
	}
	file, err := cfg.Db.CreateNewJson(r.Context(), database.CreateNewJsonParams{
		ID:       fileId,
		UserID:   userId,
		FileName: "New JSON File",
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

	// read json from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "cannot read request body", err)
		return
	}
	defer r.Body.Close()

	fmt.Println("Update JSON request body:", string(body))

	// create json file to replace on s3
	tempFile, err := os.CreateTemp("", fileMetadata.ID.String())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot create temp file", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	// copy request body to tempfile
	if _, err := tempFile.Write(body); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot write to temp file", err)
		return
	}

	// Reset the file pointer to the beginning
	if _, err := tempFile.Seek(0, 0); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot seek temp file", err)
		return
	}

	// upload file to s3
	if err := cfg.uploadFileToS3(r.Context(), userId, fileMetadata.ID, tempFile); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error uploading file to s3", err)
		return
	}

	fileContents, err := cfg.getFileFromS3(r.Context(), userId, fileMetadata.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot get updated json file from s3", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(fileContents); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error writing response", err)
		return
	}
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
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)

	fileContents, err := cfg.getFileFromS3(r.Context(), userId, fileMetadata.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot get json file from s3", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(fileContents); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error writing response", err)
		return
	}
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
			ID:       file.ID,
			UserID:   file.UserID,
			FileName: file.FileName,
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
