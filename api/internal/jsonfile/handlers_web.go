package jsonfile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"time"

	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/s3util"
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

type Route struct {
	Method      string `json:"method"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

func (cfg *JsonConfig) HandlerCreateJson(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)

	user, err := cfg.Db.GetUserById(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error getting user", err)
		return
	}
	var fileLimit int
	if user.Subscribed {
		fileLimit = cfg.ProFileLimit
	} else {
		fileLimit = cfg.FreeFileLimit
	}

	existingJsonMetadata, err := cfg.Db.GetJsonFiles(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error checking number of json files", err)
		return
	}
	if len(existingJsonMetadata) >= fileLimit {
		utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("json file limit of %d exceeded", fileLimit), err)
		return
	}

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
	err = s3util.UploadJsonToS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileId, emptyJson)
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

	if err := s3util.UploadJsonToS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID, jsonData); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error uploading JSON to s3", err)
		return
	}

	fileContents, err := s3util.GetJsonFromS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID)
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

	jsonFilesResponse := []JsonMetadataResponse{}
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
	if err := s3util.DeleteFileFromS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error deleting json file from s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *JsonConfig) HandlerGetDynamicRoutes(w http.ResponseWriter, r *http.Request) {
	fileContents, ok := r.Context().Value(FileContentContextKey).(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "json file is not a map", nil)
		return
	}
	routes := []Route{}
	for key, val := range fileContents {
		if strings.Contains(key, " ") {
			// skip keys with spaces, for example "hello world"
			// since this results in invalid url
			continue
		}
		switch val.(type) {
		case map[string]any:
			routes = append(routes, Route{
				Method:      "GET",
				Url:         fmt.Sprintf("/%s", key),
				Description: "Gets the entire resource",
			})
			routes = append(routes, Route{
				Method:      "PUT",
				Url:         fmt.Sprintf("/%s", key),
				Description: "Replaces the entire resource",
			})
			routes = append(routes, Route{
				Method:      "PATCH",
				Url:         fmt.Sprintf("/%s", key),
				Description: "Partially updates the resource",
			})
		case []any:
			routes = append(routes, Route{
				Method:      "GET",
				Url:         fmt.Sprintf("/%s", key),
				Description: "Gets the entire resource array",
			})
			routes = append(routes, Route{
				Method:      "GET",
				Url:         fmt.Sprintf("/%s/:id", key),
				Description: "Gets a resource from resource array by id",
			})
			routes = append(routes, Route{
				Method:      "POST",
				Url:         fmt.Sprintf("/%s", key),
				Description: "Creates a new resource and adds it to the resource array",
			})
			routes = append(routes, Route{
				Method:      "PUT",
				Url:         fmt.Sprintf("/%s/:id", key),
				Description: "Replaces a resource from resource array with id",
			})
			routes = append(routes, Route{
				Method:      "PATCH",
				Url:         fmt.Sprintf("/%s/:id", key),
				Description: "Partially updates a resource from resource array with id",
			})
			routes = append(routes, Route{
				Method:      "DELETE",
				Url:         fmt.Sprintf("/%s/:id", key),
				Description: "Deletes a resource from resource array with id",
			})
		case string, int, float64, bool, nil:
			routes = append(routes, Route{
				Method:      "GET",
				Url:         fmt.Sprintf("/%s", key),
				Description: "Gets the entire resource",
			})
			routes = append(routes, Route{
				Method:      "PUT",
				Url:         fmt.Sprintf("/%s", key),
				Description: "Replaces the entire resource",
			})
		}
	}
	utils.RespondWithJSON(w, http.StatusOK, routes)
}
