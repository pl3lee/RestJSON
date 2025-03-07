package jsonfile

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pl3lee/webjson/internal/auth"
	"github.com/pl3lee/webjson/internal/database"
	"github.com/pl3lee/webjson/internal/utils"
)

func (cfg *JsonConfig) HandlerCreateJson(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)

	// TODO: create file in s3
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
