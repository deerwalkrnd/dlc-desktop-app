// Package catalog turns DLC video filenames into structured lesson metadata.
//
// The filename convention is:
//
//	TEACHER - LESSON NO - CHAPTER TITLE - LESSON TITLE - SUBJECT - SYLLABUS - CLASS.mp4
//
// for example:
//
//	Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4
//
// The chapter number is the whole part of the lesson number (1.2 belongs to
// chapter 1). Parsing is a pure function of the filename and never touches the
// filesystem.
//
// # Compatibility
//
// Earlier releases split the name on every "-" character. Any filename that
// parsed under those rules parses to exactly the same result here: the legacy
// split is always tried first. Only when it fails does the parser fall back to
// splitting on " - " (a hyphen with whitespace on both sides), which lets
// fields themselves contain hyphens, e.g. "Ram-Prasad Sharma" or "Mid-term".
package catalog

import (
	"errors"
	"fmt"
	"math"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// FieldCount is the number of " - " separated fields in a valid filename.
const FieldCount = 7

// Syllabus distinguishes old and new curriculum content.
type Syllabus int

const (
	// SyllabusOld is selected by the literal field value "OLD".
	SyllabusOld Syllabus = iota
	// SyllabusNew is used for every other value (legacy behaviour).
	SyllabusNew
)

func (s Syllabus) String() string {
	if s == SyllabusOld {
		return "old"
	}
	return "new"
}

// ParseSyllabus parses the API/URL spelling ("old" or "new").
func ParseSyllabus(v string) (Syllabus, bool) {
	switch v {
	case "old":
		return SyllabusOld, true
	case "new":
		return SyllabusNew, true
	}
	return 0, false
}

// Entry is the metadata encoded in one video filename.
type Entry struct {
	Teacher       string
	LessonNumber  string // as written, e.g. "1.10"
	ChapterNumber int    // whole part of LessonNumber
	ChapterTitle  string
	LessonTitle   string
	Subject       string
	Syllabus      Syllabus
	Class         int

	// Warnings describe things that parsed but are probably mistakes.
	Warnings []string
}

// ErrNotVideo is returned for files that are not .mp4 videos.
var ErrNotVideo = errors.New("not an .mp4 file")

// spacedHyphen matches a hyphen with whitespace on both sides.
var spacedHyphen = regexp.MustCompile(`\s+-\s+`)

// IsVideo reports whether name has an .mp4 extension (any letter case).
func IsVideo(name string) bool {
	_, ok := stem(name)
	return ok
}

func stem(name string) (string, bool) {
	base := path.Base(strings.ReplaceAll(name, `\`, "/"))
	ext := path.Ext(base)
	if !strings.EqualFold(ext, ".mp4") || len(base) <= len(ext) {
		return "", false
	}
	return base[:len(base)-len(ext)], true
}

// Parse extracts lesson metadata from a filename or path. Only the base name
// is considered.
func Parse(name string) (Entry, error) {
	s, ok := stem(name)
	if !ok {
		return Entry{}, ErrNotVideo
	}

	// 1. Legacy rules: split on every hyphen.
	if e, err := fromFields(strings.Split(s, "-")); err == nil {
		return e, nil
	}

	// 2. Hyphen-tolerant rules: split on " - " only.
	fields := spacedHyphen.Split(s, -1)
	if len(fields) < FieldCount {
		return Entry{}, fmt.Errorf("expected %d parts separated by \" - \", found %d", FieldCount, len(fields))
	}
	return fromFields(fields)
}

func fromFields(items []string) (Entry, error) {
	if len(items) < FieldCount {
		return Entry{}, fmt.Errorf("expected %d parts separated by \" - \", found %d", FieldCount, len(items))
	}
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	number, err := strconv.ParseFloat(items[1], 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) || number < 0 || number >= math.MaxInt32 {
		return Entry{}, fmt.Errorf("lesson number %q is not a number like 1.2", items[1])
	}
	class, err := strconv.Atoi(items[6])
	if err != nil || class < 0 {
		return Entry{}, fmt.Errorf("class %q is not a whole number", items[6])
	}

	e := Entry{
		Teacher:       items[0],
		LessonNumber:  items[1],
		ChapterNumber: int(math.Floor(number)),
		ChapterTitle:  items[2],
		LessonTitle:   items[3],
		Subject:       items[4],
		Syllabus:      SyllabusNew,
		Class:         class,
	}

	switch items[5] {
	case "OLD":
		e.Syllabus = SyllabusOld
	case "NEW":
	default:
		hint := ""
		if strings.EqualFold(items[5], "old") {
			hint = ` (write "OLD" in capitals for the old syllabus)`
		}
		e.Warn("syllabus %q is not OLD or NEW, treated as NEW%s", items[5], hint)
	}

	for _, f := range []struct{ name, value string }{
		{"teacher", e.Teacher}, {"chapter title", e.ChapterTitle},
		{"lesson title", e.LessonTitle}, {"subject", e.Subject},
	} {
		if f.value == "" {
			e.Warn("%s is empty", f.name)
		}
	}
	if len(items) > FieldCount {
		e.Warn("%d extra part(s) after the class number were ignored", len(items)-FieldCount)
	}
	return e, nil
}

// Warn appends a formatted warning.
func (e *Entry) Warn(format string, args ...any) {
	e.Warnings = append(e.Warnings, fmt.Sprintf(format, args...))
}

// CompareLessonNumbers orders lesson numbers like versions: 1.2 < 1.10 < 2.
// Numbers that are not of the form MAJOR[.MINOR] sort after all others,
// ordered as plain strings.
func CompareLessonNumbers(a, b string) int {
	amaj, amin, aok := splitLessonNumber(a)
	bmaj, bmin, bok := splitLessonNumber(b)
	switch {
	case aok && bok:
		if amaj != bmaj {
			return cmpInt(amaj, bmaj)
		}
		return cmpInt(amin, bmin)
	case aok:
		return -1
	case bok:
		return 1
	}
	return strings.Compare(a, b)
}

func splitLessonNumber(s string) (major, minor int, ok bool) {
	majStr, minStr, hasDot := strings.Cut(s, ".")
	major, err := strconv.Atoi(majStr)
	if err != nil {
		return 0, 0, false
	}
	if !hasDot {
		return major, 0, true
	}
	minor, err = strconv.Atoi(minStr)
	if err != nil {
		return 0, 0, false
	}
	return major, minor, true
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}
