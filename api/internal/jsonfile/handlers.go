package jsonfile

import (
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

	// upload file to s3
	_, err = cfg.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3Bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", userId.String(), fileId.String())),
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
		Url:      fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/%s", cfg.S3Bucket, cfg.S3Region, userId.String(), fileId.String()),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "cannot create new json", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, file)
}
