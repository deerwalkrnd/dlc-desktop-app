package api

import (
	"net/http"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
	"gorm.io/gorm"
)

type ApiHandler struct {
	db *gorm.DB
}

func NewApiHandler(db *gorm.DB) *ApiHandler {

	return &ApiHandler{
		db: db,
	}
}

func (a *ApiHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	var teachers []db.Teacher

	result := a.db.Find(&teachers)

	if result.Error != nil {
		respondWithJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": result.Error.Error(),
			},
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		map[string]interface{}{
			"teachers": teachers,
			"count":    len(teachers),
		},
	)
}
