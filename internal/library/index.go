// Package library scans a DLC video folder and builds an immutable, in-memory
// catalogue of classes, subjects, chapters and lessons.
//
// The folder is the single source of truth: there is no database to go stale.
// A scan produces a new *Index; Library swaps it in atomically so readers
// never observe a half-built catalogue.
package library

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/deerwalkrnd/dlc-desktop-app/internal/catalog"
)

// Index is a read-only snapshot of the library. Never mutate it after Build.
type Index struct {
	Root      string
	Exists    bool // false when the library folder is missing
	ScannedAt time.Time
	Duration  time.Duration

	Classes  []*Class
	Teachers []*Teacher

	Skipped  []Issue // files that look like videos but could not be parsed
	Warnings []Issue // files that parsed, but probably contain a mistake

	classes map[int]*Class
	lessons map[string]*Lesson
	search  []searchDoc
}

// Issue ties a message to a file, relative to the library root.
type Issue struct {
	Path    string
	Message string
}

type Class struct {
	Number   int
	Subjects []*Subject
}

type Subject struct {
	Slug     string // unique within its class, stable across scans
	Name     string
	Syllabus catalog.Syllabus
	Class    int
	Chapters []*Chapter
}

type Chapter struct {
	Number  int
	Title   string
	Lessons []*Lesson
}

type Lesson struct {
	ID      string // stable hash of Path
	Number  string
	Title   string
	Teacher string
	Path    string // slash-separated, relative to the library root
}

type Teacher struct {
	Name     string
	Lessons  int
	Subjects []string // distinct subject names, sorted
}

// Class returns the class with the given number, or nil.
func (x *Index) Class(number int) *Class { return x.classes[number] }

// Lesson returns the lesson with the given ID, or nil.
func (x *Index) Lesson(id string) *Lesson { return x.lessons[id] }

// LessonCount is the number of playable lessons in the index.
func (x *Index) LessonCount() int { return len(x.lessons) }

// Subject returns the subject with the given slug, or nil.
func (c *Class) Subject(slug string) *Subject {
	for _, s := range c.Subjects {
		if s.Slug == slug {
			return s
		}
	}
	return nil
}

// LessonCount is the number of lessons in the subject.
func (s *Subject) LessonCount() int {
	n := 0
	for _, ch := range s.Chapters {
		n += len(ch.Lessons)
	}
	return n
}

// Build walks root and returns a fresh index. A missing root is not an error:
// it yields an empty index with Exists=false.
func Build(root string) (*Index, error) {
	start := time.Now()
	x := &Index{
		Root:      root,
		ScannedAt: start,
		classes:   map[int]*Class{},
		lessons:   map[string]*Lesson{},
	}

	info, err := os.Stat(root)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		x.Duration = time.Since(start)
		return x, nil
	case err != nil:
		return nil, err
	case !info.IsDir():
		return nil, fmt.Errorf("library path %s is not a folder", root)
	}
	x.Exists = true

	b := newBuilder(x)
	err = fs.WalkDir(os.DirFS(root), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if p == "." {
				return err
			}
			x.Skipped = append(x.Skipped, Issue{p, "could not read: " + err.Error()})
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() || !catalog.IsVideo(p) {
			return nil
		}
		if strings.HasPrefix(d.Name(), "._") {
			// macOS writes these metadata files next to every video copied
			// onto a FAT/exFAT USB drive. They are not videos.
			return nil
		}
		e, err := catalog.Parse(d.Name())
		if err != nil {
			x.Skipped = append(x.Skipped, Issue{p, err.Error()})
			return nil
		}
		for _, w := range e.Warnings {
			x.Warnings = append(x.Warnings, Issue{p, w})
		}
		b.add(p, e)
		return nil
	})
	if err != nil {
		return nil, err
	}
	b.finish()
	x.buildSearch()
	x.Duration = time.Since(start)
	return x, nil
}

type subjectKey struct {
	class    int
	name     string
	syllabus catalog.Syllabus
}

type chapterKey struct {
	number int
	title  string
}

type lessonKey struct {
	number, title, teacher string
}

type builder struct {
	x          *Index
	subjects   map[subjectKey]*Subject
	chapters   map[*Subject]map[chapterKey]*Chapter
	teachers   map[string]*Teacher
	teacherSub map[string]map[string]bool
	seen       map[*Chapter]map[lessonKey]string
}

func newBuilder(x *Index) *builder {
	return &builder{
		x:          x,
		subjects:   map[subjectKey]*Subject{},
		chapters:   map[*Subject]map[chapterKey]*Chapter{},
		teachers:   map[string]*Teacher{},
		teacherSub: map[string]map[string]bool{},
		seen:       map[*Chapter]map[lessonKey]string{},
	}
}

func (b *builder) add(p string, e catalog.Entry) {
	class := b.x.classes[e.Class]
	if class == nil {
		class = &Class{Number: e.Class}
		b.x.classes[e.Class] = class
		b.x.Classes = append(b.x.Classes, class)
	}

	sk := subjectKey{e.Class, e.Subject, e.Syllabus}
	subject := b.subjects[sk]
	if subject == nil {
		subject = &Subject{Name: e.Subject, Syllabus: e.Syllabus, Class: e.Class}
		b.subjects[sk] = subject
		b.chapters[subject] = map[chapterKey]*Chapter{}
		class.Subjects = append(class.Subjects, subject)
	}

	ck := chapterKey{e.ChapterNumber, e.ChapterTitle}
	chapter := b.chapters[subject][ck]
	if chapter == nil {
		chapter = &Chapter{Number: e.ChapterNumber, Title: e.ChapterTitle}
		b.chapters[subject][ck] = chapter
		b.seen[chapter] = map[lessonKey]string{}
		subject.Chapters = append(subject.Chapters, chapter)
	}

	lk := lessonKey{e.LessonNumber, e.LessonTitle, e.Teacher}
	if other, dup := b.seen[chapter][lk]; dup {
		b.x.Warnings = append(b.x.Warnings, Issue{p, "same lesson as " + other})
	} else {
		b.seen[chapter][lk] = p
	}

	lesson := &Lesson{
		ID:      lessonID(p),
		Number:  e.LessonNumber,
		Title:   e.LessonTitle,
		Teacher: e.Teacher,
		Path:    p,
	}
	chapter.Lessons = append(chapter.Lessons, lesson)
	b.x.lessons[lesson.ID] = lesson

	t := b.teachers[e.Teacher]
	if t == nil {
		t = &Teacher{Name: e.Teacher}
		b.teachers[e.Teacher] = t
		b.teacherSub[e.Teacher] = map[string]bool{}
		b.x.Teachers = append(b.x.Teachers, t)
	}
	t.Lessons++
	if !b.teacherSub[e.Teacher][e.Subject] {
		b.teacherSub[e.Teacher][e.Subject] = true
		t.Subjects = append(t.Subjects, e.Subject)
	}
}

// finish sorts everything into display order and assigns subject slugs.
func (b *builder) finish() {
	x := b.x
	slices.SortFunc(x.Classes, func(a, b *Class) int { return cmp.Compare(a.Number, b.Number) })
	for _, c := range x.Classes {
		slices.SortFunc(c.Subjects, func(a, b *Subject) int {
			return cmp.Or(strings.Compare(a.Name, b.Name), cmp.Compare(b.Syllabus, a.Syllabus))
		})
		used := map[string]bool{}
		for _, s := range c.Subjects {
			s.Slug = uniqueSlug(slugify(s.Name)+"-"+s.Syllabus.String(), used)
			slices.SortFunc(s.Chapters, func(a, b *Chapter) int {
				return cmp.Or(cmp.Compare(a.Number, b.Number), strings.Compare(a.Title, b.Title))
			})
			for _, ch := range s.Chapters {
				slices.SortFunc(ch.Lessons, func(a, b *Lesson) int {
					return cmp.Or(
						catalog.CompareLessonNumbers(a.Number, b.Number),
						strings.Compare(a.Title, b.Title),
						strings.Compare(a.Path, b.Path),
					)
				})
			}
		}
	}
	slices.SortFunc(x.Teachers, func(a, b *Teacher) int { return strings.Compare(a.Name, b.Name) })
	for _, t := range x.Teachers {
		slices.Sort(t.Subjects)
	}
	byPath := func(a, b Issue) int {
		return cmp.Or(strings.Compare(a.Path, b.Path), strings.Compare(a.Message, b.Message))
	}
	slices.SortFunc(x.Skipped, byPath)
	slices.SortFunc(x.Warnings, byPath)
}

// lessonID derives a short, stable, URL-safe identifier from a file path so
// that links keep working across rescans and restarts.
func lessonID(p string) string {
	sum := sha256.Sum256([]byte(p))
	return hex.EncodeToString(sum[:8])
}

// slugify lower-cases s and joins runs of letters, marks and digits with "-".
// Non-Latin scripts (e.g. Devanagari) are kept as-is.
func slugify(s string) string {
	var sb strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r) {
			if dash && sb.Len() > 0 {
				sb.WriteByte('-')
			}
			sb.WriteRune(r)
			dash = false
		} else {
			dash = true
		}
	}
	if sb.Len() == 0 {
		return "subject"
	}
	return sb.String()
}

func uniqueSlug(base string, used map[string]bool) string {
	slug := base
	for i := 2; used[slug]; i++ {
		slug = base + "-" + strconv.Itoa(i)
	}
	used[slug] = true
	return slug
}
