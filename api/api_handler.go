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
	router.HandleFunc("/api/teachers", a.GetTeachers).Methods("GET")
	router.HandleFunc("/api/classes", a.GetClasses).Methods("GET")
	router.HandleFunc("/api/classes/{classID}/subjects", a.GetSubjectsByClass).Methods("GET")
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

func (a *ApiHandler) GetSubjectsByClass(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	vars := mux.Vars(r)
	classID := vars["classID"]

	typeParam := r.URL.Query().Get("type")
	type_ := db.NEW_SYLLABUS

	if typeParam == "" || typeParam == "old" {
		type_ = db.OLD_SYLLABUS
	}

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

	var subjects []db.Subject

	query := a.db.Where("class_id = ?", classID)

	if typeParam != "" {
		query = query.Where("type = ?", type_)
	}

	result := query.Find(&subjects)

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
			"subjects": subjects,
			"count":    len(subjects),
		},
	)
}
