package api

import (
	"net/http"
	"strconv"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
	"github.com/gorilla/mux"
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

func (a *ApiHandler) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/teachers", a.GetTeachers).Methods("GET")
	router.HandleFunc("/classes", a.GetClasses).Methods("GET")
	router.HandleFunc("/classes/{classID}/lectures", a.GetLecturesByClass).Methods("GET")
}

func (a *ApiHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	var teachers []db.Teacher

	result := a.db.Order("name asc").Find(&teachers)

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
	enableCORS(w)
	var classes []db.Class

	result := a.db.Order("number asc").Find(&classes)

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

func (a *ApiHandler) GetLecturesByClass(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	vars := mux.Vars(r)
	classID := vars["classID"]

	typeParam := r.URL.Query().Get("type")

	_, err := strconv.Atoi(classID)
	if err != nil {
		respondWithJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid class ID",
			},
		)
		return
	}

	var lectures []db.Lecture
	query := a.db.Where("class_id = ?", classID)

	if typeParam != "" {
		query = query.Where("type = ?", typeParam)
	}

	result := query.Find(&lectures)

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
			"lectures": lectures,
			"count":    len(lectures),
		},
	)
}
