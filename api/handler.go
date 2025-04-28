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

	result := a.db.Find(&teachers).Order("name asc")

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
func (a *ApiHandler) GetClasses(w http.ResponseWriter, r *http.Request) {
	var classes []db.Class

	result := a.db.Find(&classes).Statement.Order("number asc")

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
			"classes": classes,
			"count":   len(classes),
		},
	)
}
