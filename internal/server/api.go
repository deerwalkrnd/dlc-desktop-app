package server

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/deerwalkrnd/dlc-desktop-app/internal/catalog"
	"github.com/deerwalkrnd/dlc-desktop-app/internal/library"
)

// JSON response shapes. Keep in sync with web/src/lib/api.ts.

type classSummary struct {
	Number       int      `json:"number"`
	SubjectCount int      `json:"subjectCount"`
	LessonCount  int      `json:"lessonCount"`
	Syllabi      []string `json:"syllabi"`
}

type subjectSummary struct {
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Syllabus     string `json:"syllabus"`
	ChapterCount int    `json:"chapterCount"`
	LessonCount  int    `json:"lessonCount"`
}

type classDetail struct {
	Number   int              `json:"number"`
	Syllabi  []string         `json:"syllabi"`
	Subjects []subjectSummary `json:"subjects"`
}

type subjectDetail struct {
	Class    int          `json:"class"`
	Slug     string       `json:"slug"`
	Name     string       `json:"name"`
	Syllabus string       `json:"syllabus"`
	Chapters []chapterDTO `json:"chapters"`
}

type chapterDTO struct {
	Number  int         `json:"number"`
	Title   string      `json:"title"`
	Lessons []lessonDTO `json:"lessons"`
}

type lessonDTO struct {
	ID       string `json:"id"`
	Number   string `json:"number"`
	Title    string `json:"title"`
	Teacher  string `json:"teacher"`
	MediaURL string `json:"mediaUrl"`
}

type teacherDTO struct {
	Name        string   `json:"name"`
	LessonCount int      `json:"lessonCount"`
	Subjects    []string `json:"subjects"`
}

type searchHitDTO struct {
	Class         int       `json:"class"`
	SubjectSlug   string    `json:"subjectSlug"`
	SubjectName   string    `json:"subjectName"`
	Syllabus      string    `json:"syllabus"`
	ChapterNumber int       `json:"chapterNumber"`
	ChapterTitle  string    `json:"chapterTitle"`
	Lesson        lessonDTO `json:"lesson"`
}

type searchDTO struct {
	Query   string         `json:"query"`
	Total   int            `json:"total"`
	Results []searchHitDTO `json:"results"`
}

type issueDTO struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type libraryDTO struct {
	Path        string     `json:"path"`
	Exists      bool       `json:"exists"`
	ScannedAt   time.Time  `json:"scannedAt"`
	ScanMillis  int64      `json:"scanMillis"`
	ClassCount  int        `json:"classCount"`
	LessonCount int        `json:"lessonCount"`
	Skipped     []issueDTO `json:"skipped"`
	Warnings    []issueDTO `json:"warnings"`
	Version     string     `json:"version"`
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "app": "dlc-desktop-app", "version": s.version})
}

func (s *Server) listClasses(w http.ResponseWriter, r *http.Request) {
	x := s.lib.Index()
	out := make([]classSummary, 0, len(x.Classes))
	for _, c := range x.Classes {
		cs := classSummary{Number: c.Number, SubjectCount: len(c.Subjects), Syllabi: syllabi(c)}
		for _, sub := range c.Subjects {
			cs.LessonCount += sub.LessonCount()
		}
		out = append(out, cs)
	}
	writeJSON(w, http.StatusOK, map[string]any{"classes": out})
}

func (s *Server) getClass(w http.ResponseWriter, r *http.Request) {
	c, ok := s.lookupClass(w, r)
	if !ok {
		return
	}
	filter, ok := parseSyllabusFilter(w, r)
	if !ok {
		return
	}
	out := classDetail{Number: c.Number, Syllabi: syllabi(c), Subjects: []subjectSummary{}}
	for _, sub := range c.Subjects {
		if filter != nil && sub.Syllabus != *filter {
			continue
		}
		out.Subjects = append(out.Subjects, subjectSummary{
			Slug:         sub.Slug,
			Name:         sub.Name,
			Syllabus:     sub.Syllabus.String(),
			ChapterCount: len(sub.Chapters),
			LessonCount:  sub.LessonCount(),
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) getSubject(w http.ResponseWriter, r *http.Request) {
	c, ok := s.lookupClass(w, r)
	if !ok {
		return
	}
	sub := c.Subject(r.PathValue("slug"))
	if sub == nil {
		writeError(w, http.StatusNotFound, "subject not found")
		return
	}
	out := subjectDetail{
		Class:    c.Number,
		Slug:     sub.Slug,
		Name:     sub.Name,
		Syllabus: sub.Syllabus.String(),
		Chapters: make([]chapterDTO, 0, len(sub.Chapters)),
	}
	for _, ch := range sub.Chapters {
		cd := chapterDTO{Number: ch.Number, Title: ch.Title, Lessons: make([]lessonDTO, 0, len(ch.Lessons))}
		for _, l := range ch.Lessons {
			cd.Lessons = append(cd.Lessons, toLessonDTO(l))
		}
		out.Chapters = append(out.Chapters, cd)
	}
	writeJSON(w, http.StatusOK, out)
}

const (
	defaultSearchLimit = 50
	maxSearchLimit     = 200
	maxQueryLength     = 200
)

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) > maxQueryLength {
		writeError(w, http.StatusBadRequest, "search query is too long")
		return
	}
	limit := defaultSearchLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			writeError(w, http.StatusBadRequest, "limit must be a positive number")
			return
		}
		limit = min(n, maxSearchLimit)
	}

	hits, total := s.lib.Index().Search(q, limit)
	out := searchDTO{Query: q, Total: total, Results: make([]searchHitDTO, 0, len(hits))}
	for _, h := range hits {
		out.Results = append(out.Results, searchHitDTO{
			Class:         h.Subject.Class,
			SubjectSlug:   h.Subject.Slug,
			SubjectName:   h.Subject.Name,
			Syllabus:      h.Subject.Syllabus.String(),
			ChapterNumber: h.Chapter.Number,
			ChapterTitle:  h.Chapter.Title,
			Lesson:        toLessonDTO(h.Lesson),
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func toLessonDTO(l *library.Lesson) lessonDTO {
	return lessonDTO{
		ID:       l.ID,
		Number:   l.Number,
		Title:    l.Title,
		Teacher:  l.Teacher,
		MediaURL: "/media/" + l.ID,
	}
}

func (s *Server) listTeachers(w http.ResponseWriter, r *http.Request) {
	x := s.lib.Index()
	out := make([]teacherDTO, 0, len(x.Teachers))
	for _, t := range x.Teachers {
		out = append(out, teacherDTO{Name: t.Name, LessonCount: t.Lessons, Subjects: t.Subjects})
	}
	writeJSON(w, http.StatusOK, map[string]any{"teachers": out})
}

func (s *Server) libraryStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.libraryDTO(s.lib.Index()))
}

func (s *Server) rescan(w http.ResponseWriter, r *http.Request) {
	if !sameOrigin(r) {
		writeError(w, http.StatusForbidden, "cross-origin request rejected")
		return
	}
	x, err := s.lib.Rescan()
	if err != nil {
		s.log.Error("rescan failed", "err", err)
		writeError(w, http.StatusInternalServerError, "could not scan the library folder: "+err.Error())
		return
	}
	s.log.Info("library rescanned", "lessons", x.LessonCount(), "skipped", len(x.Skipped), "took", x.Duration)
	writeJSON(w, http.StatusOK, s.libraryDTO(x))
}

func (s *Server) libraryDTO(x *library.Index) libraryDTO {
	return libraryDTO{
		Path:        x.Root,
		Exists:      x.Exists,
		ScannedAt:   x.ScannedAt,
		ScanMillis:  x.Duration.Milliseconds(),
		ClassCount:  len(x.Classes),
		LessonCount: x.LessonCount(),
		Skipped:     issues(x.Skipped),
		Warnings:    issues(x.Warnings),
		Version:     s.version,
	}
}

func (s *Server) media(w http.ResponseWriter, r *http.Request) {
	lesson := s.lib.Index().Lesson(r.PathValue("id"))
	if lesson == nil {
		http.Error(w, "video not found", http.StatusNotFound)
		return
	}
	f, info, err := s.lib.Open(lesson)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "video file is missing; rescan the library", http.StatusNotFound)
			return
		}
		s.log.Warn("cannot open video", "path", lesson.Path, "err", err)
		http.Error(w, "cannot open video", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "private, no-cache")
	// ServeContent implements Range requests, which video seeking relies on.
	http.ServeContent(w, r, "", info.ModTime(), f)
}

func (s *Server) lookupClass(w http.ResponseWriter, r *http.Request) (*library.Class, bool) {
	n, err := strconv.Atoi(r.PathValue("number"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "class must be a number")
		return nil, false
	}
	c := s.lib.Index().Class(n)
	if c == nil {
		writeError(w, http.StatusNotFound, "class not found")
		return nil, false
	}
	return c, true
}

// parseSyllabusFilter reads ?syllabus=old|new. Absent means no filter.
func parseSyllabusFilter(w http.ResponseWriter, r *http.Request) (*catalog.Syllabus, bool) {
	v := r.URL.Query().Get("syllabus")
	if v == "" {
		return nil, true
	}
	syl, ok := catalog.ParseSyllabus(v)
	if !ok {
		writeError(w, http.StatusBadRequest, `syllabus must be "old" or "new"`)
		return nil, false
	}
	return &syl, true
}

func syllabi(c *library.Class) []string {
	var hasNew, hasOld bool
	for _, s := range c.Subjects {
		if s.Syllabus == catalog.SyllabusOld {
			hasOld = true
		} else {
			hasNew = true
		}
	}
	out := []string{}
	if hasNew {
		out = append(out, catalog.SyllabusNew.String())
	}
	if hasOld {
		out = append(out, catalog.SyllabusOld.String())
	}
	return out
}

func issues(in []library.Issue) []issueDTO {
	out := make([]issueDTO, len(in))
	for i, is := range in {
		out[i] = issueDTO{Path: is.Path, Message: is.Message}
	}
	return out
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	w.Write(body)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
