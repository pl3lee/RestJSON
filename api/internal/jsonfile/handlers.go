package jsonfile

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/pl3lee/webjson/internal/auth"
	"github.com/pl3lee/webjson/internal/database"
	"github.com/pl3lee/webjson/internal/utils"
)

func (cfg *JsonConfig) HandlerCreateJson(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)

	// read json from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "cannot read request body", err)
		return
	}
	defer r.Body.Close()

	fmt.Println("Request Body:", string(body))

	// create json file
	fileId := uuid.New()
	tempFile, err := os.CreateTemp("", fileId.String())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot create temp file", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
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
	_, err = cfg.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3Bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s.json", userId.String(), fileId.String())),
		Body:        tempFile,
		ContentType: aws.String("application/json"),
	})
	if err != nil {
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

func (cfg *JsonConfig) HandlerGetJson(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileId := r.Context().Value(FileIDContextKey).(uuid.UUID)

	s3Params := &s3.GetObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s.json", userId.String(), fileId.String())),
	}

	jsonFile, err := cfg.S3Client.GetObject(r.Context(), s3Params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot get json file from s3", err)
		return
	}
	defer jsonFile.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// copy the content from the S3 object to the response writer
	if _, err := io.Copy(w, jsonFile.Body); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error writing json file to response", err)
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
	utils.RespondWithJSON(w, http.StatusOK, jsonFiles)
}

type RenameRequest struct {
	FileName string `json:"fileName"`
}

func (cfg *JsonConfig) HandlerRenameJsonFile(w http.ResponseWriter, r *http.Request) {
	fileId := r.Context().Value(FileIDContextKey).(uuid.UUID)

	var renameReq RenameRequest
	if err := json.NewDecoder(r.Body).Decode(&renameReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	defer r.Body.Close()

	if renameReq.FileName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "file name cannot be empty", nil)
		return
	}

	newJsonFile, err := cfg.Db.RenameJsonFile(r.Context(), database.RenameJsonFileParams{
		ID:       fileId,
		FileName: renameReq.FileName,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error renaming json file", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, newJsonFile)
}
