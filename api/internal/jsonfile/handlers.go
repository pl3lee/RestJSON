package jsonfile

import (
	"io"
	"net/http"
	"os"

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
	tempFile, err := os.CreateTemp("", "temp.json")
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

	// TODO: upload file to s3
	file, err := cfg.Db.CreateNewJson(r.Context(), database.CreateNewJsonParams{
		UserID:   userId,
		FileName: "file",
		Url:      "",
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "cannot create new json", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, file)
}
