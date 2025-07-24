package api

import (
	"net/http"
	"strconv"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SimplifiedChapter struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Number    uint   `json:"number"`
	VideoUrl  string `json:"videoUrl"`
	TeacherID uint   `json:"teacherId"`
	Teacher   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"teacher,omitempty"`
}

type SimplifiedUnit struct {
	ID        uint                `json:"id"`
	Number    uint                `json:"number"`
	Name      string              `json:"name"`
	SubjectID uint                `json:"subjectId"`
	Chapters  []SimplifiedChapter `json:"chapters,omitempty"`
}

type ApiHandler struct {
	db *gorm.DB
}

func NewApiHandler(db *gorm.DB) *ApiHandler {
	return &ApiHandler{db: db}
}

func (a *ApiHandler) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/api/teachers", a.GetTeachers).Methods("GET")
	router.HandleFunc("/api/classes", a.GetClasses).Methods("GET")
	router.HandleFunc("/api/classes/{classID}/subjects", a.GetSubjectsByClass).Methods("GET")
	router.HandleFunc("/api/subjects/{subjectId}/units", a.GetUnitsBySubject).Methods("GET")
	router.HandleFunc("/api/units/{unitId}/chapters", a.GetChaptersByUnit).Methods("GET")
}

func (a *ApiHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	var teachers []db.Teacher
	if err := a.db.Order("name asc").Find(&teachers).Error; err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]any{"teachers": teachers, "count": len(teachers)})
}

func (a *ApiHandler) GetClasses(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	var classes []db.Class
	if err := a.db.Order("number asc").Find(&classes).Error; err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]any{"classes": classes, "count": len(classes)})
}

func (a *ApiHandler) GetSubjectsByClass(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	classID := mux.Vars(r)["classID"]

	if _, err := strconv.Atoi(classID); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid class ID"})
		return
	}

	typeParam := r.URL.Query().Get("type")
	subjectType := db.NEW_SYLLABUS
	if typeParam == "" || typeParam == "old" {
		subjectType = db.OLD_SYLLABUS
	}

	var subjects []db.Subject
	query := a.db.Where("class_id = ?", classID)
	if typeParam != "" {
		query = query.Where("type = ?", subjectType)
	}
	if err := query.Find(&subjects).Error; err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"subjects": subjects, "count": len(subjects)})
}

func (a *ApiHandler) GetUnitsBySubject(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	subjectId := mux.Vars(r)["subjectId"]

	if _, err := strconv.Atoi(subjectId); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid subject ID"})
		return
	}

	var units []db.Unit
	if err := a.db.
		Where("subject_id = ?", subjectId).
		Preload("Chapters.Teacher").
		Order("number asc").
		Find(&units).Error; err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	simplified := make([]SimplifiedUnit, len(units))
	for i, u := range units {
		simplified[i] = SimplifiedUnit{
			ID:        u.ID,
			Number:    u.Number,
			Name:      u.Name,
			SubjectID: u.SubjectId,
			Chapters:  make([]SimplifiedChapter, len(u.Chapters)),
		}
		for j, c := range u.Chapters {
			simplified[i].Chapters[j] = SimplifiedChapter{
				ID:        c.ID,
				Name:      c.Name,
				Number:    c.Number,
				VideoUrl:  c.VideoUrl,
				TeacherID: c.TeacherId,
			}
			if c.Teacher.ID != 0 {
				simplified[i].Chapters[j].Teacher.ID = c.Teacher.ID
				simplified[i].Chapters[j].Teacher.Name = c.Teacher.Name
			}
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"units": simplified, "count": len(simplified)})
}

func (a *ApiHandler) GetChaptersByUnit(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	unitId := mux.Vars(r)["unitId"]

	if _, err := strconv.Atoi(unitId); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid unit ID"})
		return
	}

	var chapters []db.Chapter
	if err := a.db.
		Where("unit_id = ?", unitId).
		Preload("Teacher").
		Order("number asc").
		Find(&chapters).Error; err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := make([]SimplifiedChapter, len(chapters))
	for i, c := range chapters {
		resp[i] = SimplifiedChapter{
			ID:        c.ID,
			Name:      c.Name,
			Number:    c.Number,
			VideoUrl:  c.VideoUrl,
			TeacherID: c.TeacherId,
		}
		if c.Teacher.ID != 0 {
			resp[i].Teacher.ID = c.Teacher.ID
			resp[i].Teacher.Name = c.Teacher.Name
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"chapters": resp, "count": len(resp)})
}
