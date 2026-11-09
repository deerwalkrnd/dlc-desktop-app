package catalog

import (
	"math"
	"path/filepath"
	"strconv"
	"strings"
)

// legacyVideo and legacyParse are a verbatim copy (minus logging, plus the
// RawSubjectType field used to model the deliberate syllabus change) of
// data.ParseVideoV2 from the SQLite-based releases. They are the oracle for the
// compatibility guarantee: every filename the old app accepted must parse to
// the same metadata with the new parser.
type legacyVideo struct {
	Class         uint
	TeacherName   string
	SubjectName   string
	SubjectType   int // 0 = OLD, 1 = NEW
	LectureNumber uint
	LectureName   string
	LessonNumber  string
	LessionName   string

	RawSubjectType string // trimmed syllabus field, for test bookkeeping only
}

func legacyParse(path string) *legacyVideo {
	basePath := filepath.Base(path)

	if len(basePath) <= 1 {
		return nil
	}

	videoPath := basePath[0 : len(basePath)-4]

	items := strings.Split(videoPath, "-")

	if len(items) < 7 {
		return nil
	}
	lessionNUmberStr := strings.TrimSpace(items[1])

	lessionNumber, err := strconv.ParseFloat(strings.TrimSpace(items[1]), 64)
	if err != nil {
		return nil
	}

	classNumber, err := strconv.Atoi(strings.TrimSpace(items[6]))
	if err != nil {
		return nil
	}

	teacherName := strings.TrimSpace(items[0])
	lectureTitle := strings.TrimSpace(items[2])
	lessionTitle := strings.TrimSpace(items[3])
	subjectName := strings.TrimSpace(items[4])
	subjectTypeParsed := strings.TrimSpace(items[5])
	subjectType := 1

	if subjectTypeParsed == "OLD" {
		subjectType = 0
	}

	lectureNumber := math.Floor(lessionNumber)

	return &legacyVideo{
		TeacherName:   teacherName,
		LessonNumber:  lessionNUmberStr,
		LessionName:   lessionTitle,
		LectureName:   lectureTitle,
		SubjectName:   subjectName,
		LectureNumber: uint(lectureNumber),
		SubjectType:   subjectType,
		Class:         uint(classNumber),

		RawSubjectType: subjectTypeParsed,
	}
}

// legacyNumberSane reports whether the old parser's lesson number was a real
// number. The old parser accepted "NaN", "Inf" and negative values and turned
// them into garbage chapter numbers; the new parser rejects them on purpose.
func legacyNumberSane(s string) bool {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) && f >= 0 && f < math.MaxInt32
}

// expectedSyllabus is the legacy syllabus adjusted for the one deliberate
// change: "old" in any letter case now means the old syllabus.
func expectedSyllabus(l *legacyVideo) int {
	if strings.EqualFold(l.RawSubjectType, "OLD") {
		return 0
	}
	return l.SubjectType
}

func sameAsLegacy(e Entry, l *legacyVideo) bool {
	return uint(e.Class) == l.Class &&
		e.Teacher == l.TeacherName &&
		e.Subject == l.SubjectName &&
		int(e.Syllabus) == expectedSyllabus(l) &&
		uint(e.ChapterNumber) == l.LectureNumber &&
		e.ChapterTitle == l.LectureName &&
		e.LessonNumber == l.LessonNumber &&
		e.LessonTitle == l.LessionName
}
