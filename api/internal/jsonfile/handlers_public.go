package jsonfile

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/s3util"
	"github.com/pl3lee/restjson/internal/utils"
)

func (cfg *JsonConfig) HandlerGetResource(w http.ResponseWriter, r *http.Request) {
	resource := chi.URLParam(r, "resource")
	fileContents, ok := r.Context().Value(FileContentContextKey).(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "json file is not a map", nil)
		return
	}
	resourceData, ok := fileContents[resource]
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "resource does not exist in json", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, resourceData)
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
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
	resource := chi.URLParam(r, "resource")
	fileContents, ok := r.Context().Value(FileContentContextKey).(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "json file is not a map", nil)
		return
	}
	items, ok := r.Context().Value(ResourceArrayContextKey).([]any)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "resource is not a slice", nil)
		return
	}
	var newResource map[string]any
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&newResource); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	items = append(items, newResource)
	fileContents[resource] = items

	err := s3util.UploadJsonToS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID, fileContents)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save updated file contents to s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, fileContents)
}

func (cfg *JsonConfig) HandlerUpdateResourceItem(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
	resource := chi.URLParam(r, "resource")
	fileContents, ok := r.Context().Value(FileContentContextKey).(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "json file is not a map", nil)
		return
	}
	items, ok := r.Context().Value(ResourceArrayContextKey).([]any)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "resource is not a slice", nil)
		return
	}
	var updatedResourceItem map[string]any
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updatedResourceItem); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	resourceId := chi.URLParam(r, "id")

	foundResourceItem := false
	for index, item := range items {
		itemMap, ok := item.(map[string]any)
		if ok {
			id := fmt.Sprintf("%v", itemMap["id"])
			if id == resourceId {
				items[index] = updatedResourceItem
				foundResourceItem = true
			}
		}
	}
	if !foundResourceItem {
		utils.RespondWithError(w, http.StatusNotFound, "cannot find resource item with given id", nil)
		return
	}
	fileContents[resource] = items

	err := s3util.UploadJsonToS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID, fileContents)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save updated file contents to s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, fileContents)
}

func (cfg *JsonConfig) HandlerPartialUpdateResourceItem(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
	resource := chi.URLParam(r, "resource")
	fileContents, ok := r.Context().Value(FileContentContextKey).(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "json file is not a map", nil)
		return
	}
	items, ok := r.Context().Value(ResourceArrayContextKey).([]any)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "resource is not a slice", nil)
		return
	}
	var partialUpdate map[string]any
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&partialUpdate); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	resourceId := chi.URLParam(r, "id")

	foundResourceItem := false
	for index, item := range items {
		itemMap, ok := item.(map[string]any)
		if ok {
			id := fmt.Sprintf("%v", itemMap["id"])
			if id == resourceId {
				maps.Copy(itemMap, partialUpdate)
				items[index] = itemMap
				foundResourceItem = true
			}
		}
	}
	if !foundResourceItem {
		utils.RespondWithError(w, http.StatusNotFound, "cannot find resource item with given id", nil)
		return
	}
	fileContents[resource] = items

	err := s3util.UploadJsonToS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID, fileContents)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save updated file contents to s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, fileContents)
}

func (cfg *JsonConfig) HandlerDeleteResourceItem(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
	resource := chi.URLParam(r, "resource")
	fileContents, ok := r.Context().Value(FileContentContextKey).(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "json file is not a map", nil)
		return
	}
	items, ok := r.Context().Value(ResourceArrayContextKey).([]any)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "resource is not a slice", nil)
		return
	}
	resourceId := chi.URLParam(r, "id")

	foundResourceItem := false
	for index, item := range items {
		itemMap, ok := item.(map[string]any)
		if ok {
			id := fmt.Sprintf("%v", itemMap["id"])
			if id == resourceId {
				items = slices.Delete(items, index, index+1)
				foundResourceItem = true
				break
			}
		}
	}
	if !foundResourceItem {
		utils.RespondWithError(w, http.StatusNotFound, "cannot find resource item with given id", nil)
		return
	}
	fileContents[resource] = items

	err := s3util.UploadJsonToS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID, fileContents)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save updated file contents to s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, fileContents)
}

func (cfg *JsonConfig) HandlerUpdateResource(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
	resource := chi.URLParam(r, "resource")
	fileContents, ok := r.Context().Value(FileContentContextKey).(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "json file is not a map", nil)
		return
	}
	var updatedResource any
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updatedResource); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	fileContents[resource] = updatedResource

	err := s3util.UploadJsonToS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID, fileContents)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save updated file contents to s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, fileContents)
}

func (cfg *JsonConfig) HandlerPartialUpdateResource(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	fileMetadata := r.Context().Value(FileMetadataContextKey).(database.JsonFile)
	resource := chi.URLParam(r, "resource")
	fileContents, ok := r.Context().Value(FileContentContextKey).(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "json file is not a map", nil)
		return
	}
	existingResource, ok := fileContents[resource].(map[string]any)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "resource is not a map", nil)
		return
	}
	var partialUpdate map[string]any
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&partialUpdate); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	maps.Copy(existingResource, partialUpdate)

	fileContents[resource] = existingResource

	err := s3util.UploadJsonToS3(r.Context(), cfg.S3Client, cfg.Rdb, cfg.S3Bucket, userId, fileMetadata.ID, fileContents)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save updated file contents to s3", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, fileContents)
}
