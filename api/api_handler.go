package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SimplifiedLesson struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Number    float64   `json:"number"`
	VideoUrl  string    `json:"videoUrl"`
	TeacherID uint      `json:"teacherId"`
	LectureID uint      `json:"lectureId"`
	Teacher   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"teacher,omitempty"`
}

type SimplifiedLecture struct {
	ID        uint               `json:"id"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Number    uint               `json:"number"`
	Name      string             `json:"name"`
	SubjectID uint               `json:"subjectId"`
	Lessons   []SimplifiedLesson `json:"lessons,omitempty"`
}

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
	router.HandleFunc("/api/subjects/{subjectId}/lectures", a.GetLecturesBySubject).Methods("GET")
	router.HandleFunc("/api/lectures/{lectureId}/lessons", a.GetLessonsByLecture).Methods("GET")
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
		map[string]any{
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
		map[string]any{
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

func (a *ApiHandler) GetLecturesBySubject(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	vars := mux.Vars(r)
	subjectId := vars["subjectId"]
	_, err := strconv.Atoi(subjectId)

	if err != nil {
		respondWithJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid subject ID",
			},
		)
		return
	}

	var lectures []db.Lecture

	query := a.db.Where("subject_id = ?", subjectId)
	query = query.Preload("Lessons.Teacher")

	query = query.Preload("Lessons", func(db *gorm.DB) *gorm.DB {
		return db.Order("number asc")
	})

	result := query.Order("number asc").Find(&lectures)

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

	simplifiedLectures := make([]SimplifiedLecture, len(lectures))
	for i, lecture := range lectures {
		simplifiedLectures[i] = SimplifiedLecture{
			ID:        lecture.ID,
			CreatedAt: lecture.CreatedAt,
			UpdatedAt: lecture.UpdatedAt,
			Number:    lecture.Number,
			Name:      lecture.Name,
			SubjectID: lecture.SubjectId,
			Lessons:   make([]SimplifiedLesson, len(lecture.Lessons)),
		}

		for j, lesson := range lecture.Lessons {
			simplifiedLectures[i].Lessons[j] = SimplifiedLesson{
				ID:        lesson.ID,
				CreatedAt: lesson.CreatedAt,
				UpdatedAt: lesson.UpdatedAt,
				Name:      lesson.Name,
				Number:    lesson.Number,
				VideoUrl:  lesson.VideoUrl,
				TeacherID: lesson.TeacherId,
				LectureID: lesson.LectureId,
			}

			if lesson.Teacher.ID != 0 {
				simplifiedLectures[i].Lessons[j].Teacher.ID = lesson.Teacher.ID
				simplifiedLectures[i].Lessons[j].Teacher.Name = lesson.Teacher.Name
			}
		}
	}

	respondWithJSON(
		w,
		http.StatusOK,
		map[string]any{
			"lectures": simplifiedLectures,
			"count":    len(simplifiedLectures),
		},
	)
}

func (a *ApiHandler) GetLessonsByLecture(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	vars := mux.Vars(r)
	lectureId := vars["lectureId"]
	_, err := strconv.Atoi(lectureId)

	if err != nil {
		respondWithJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid lecture ID",
			},
		)
		return
	}

	var lessons []db.Lesson
	query := a.db.Where("lecture_id = ?", lectureId)

	result := query.Joins("Teacher").Find(&lessons)

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
		map[string]any{
			"lessons": lessons,
			"count":   len(lessons),
		},
	)

}
