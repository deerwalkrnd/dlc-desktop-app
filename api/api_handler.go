package api

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	Number    string    `json:"number"`
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

// Helper function to parse lesson number for version-style sorting
func parseLessonNumberForSorting(numberStr string) (int, int, error) {
	parts := strings.Split(numberStr, ".")
	if len(parts) == 1 {
		// No decimal point, treat as X.0
		major, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, err
		}
		return major, 0, nil
	} else if len(parts) == 2 {
		// Major.Minor format
		major, err1 := strconv.Atoi(parts[0])
		minor, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return 0, 0, err1
		}
		return major, minor, nil
	}
	return 0, 0, strconv.ErrSyntax
}

// Helper function to sort lessons by version number (1.1 < 1.2 < 1.3 < 1.10)
func sortLessons(lessons []db.Lesson) {
	sort.Slice(lessons, func(i, j int) bool {
		majorI, minorI, errI := parseLessonNumberForSorting(lessons[i].Number)
		majorJ, minorJ, errJ := parseLessonNumberForSorting(lessons[j].Number)

		// If parsing fails, fall back to string comparison
		if errI != nil || errJ != nil {
			return lessons[i].Number < lessons[j].Number
		}

		// Compare major numbers first
		if majorI != majorJ {
			return majorI < majorJ
		}

		// If major numbers are equal, compare minor numbers
		return minorI < minorJ
	})
}

// Helper function to sort simplified lessons by version number
func sortSimplifiedLessons(lessons []SimplifiedLesson) {
	sort.Slice(lessons, func(i, j int) bool {
		majorI, minorI, errI := parseLessonNumberForSorting(lessons[i].Number)
		majorJ, minorJ, errJ := parseLessonNumberForSorting(lessons[j].Number)

		// If parsing fails, fall back to string comparison
		if errI != nil || errJ != nil {
			return lessons[i].Number < lessons[j].Number
		}

		// Compare major numbers first
		if majorI != majorJ {
			return majorI < majorJ
		}

		// If major numbers are equal, compare minor numbers
		return minorI < minorJ
	})
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

	// Remove the ORDER BY from the preload since we'll sort manually
	query = query.Preload("Lessons")

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
		// Sort lessons manually for proper numerical ordering
		sortLessons(lecture.Lessons)

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

	// Sort lessons manually for proper numerical ordering
	sortLessons(lessons)

	respondWithJSON(
		w,
		http.StatusOK,
		map[string]any{
			"lessons": lessons,
			"count":   len(lessons),
		},
	)
}
