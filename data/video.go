package data

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
)

// TEACHER NAME - LESSION NO - LESSION TITLE - LECTURE NO - LECTURE TITLE - SUBJECT NAME - SUBJECT TYPE - CLASS .mp4
const MATCH_VIDEO_PATTERN = `^([^-]+) - ([^-]+) - ([^-]+) - ([^-]+) - ([^-]+) - ([^-]+) - ([^-]+) - (\d+) \.mp4$`

type Video struct {
	Class uint

	TeacherName string

	SubjectName string
	SubjectType db.SubjectType

	LectureNumber uint
	LectureName   string

	LessonNumber float64
	LessionName  string

	VideoURL string
}

func NewVideo(
	TeacherName, LessonName, LectureName, SubjectName string,
	LessionNumber float64, LectureNumber, Class uint,
	SubjectType db.SubjectType, VideoURL string) *Video {
	return &Video{
		TeacherName:   TeacherName,
		LessonNumber:  LessionNumber,
		LessionName:   LectureName,
		LectureNumber: LectureNumber,
		LectureName:   LectureName,
		SubjectName:   SubjectName,
		SubjectType:   SubjectType,
		VideoURL:      VideoURL,
		Class:         Class,
	}
}

func ParseVideo(path string) *Video {

	basePath := filepath.Base(path)
	fmt.Println("parsing: ", basePath)

	re := regexp.MustCompile(MATCH_VIDEO_PATTERN)
	matches := re.FindStringSubmatch(basePath)

	if matches == nil || len(matches) != 9 {
		log.Printf("Failed to parse, not in correct format: %s\n", path)
		return nil
	}

	lessionNumber, err := strconv.ParseFloat(strings.TrimSpace(matches[2]), 64)

	if err != nil {
		log.Printf("Could not parse lessionNumber to float64: %s\n", err.Error())
		return nil
	}

	lectureNumber, err := strconv.Atoi(strings.TrimSpace(matches[4]))

	if err != nil {
		log.Printf("Could not parse lectureNumber to int: %s\n", err.Error())
		return nil
	}

	classNumber, err := strconv.Atoi(strings.TrimSpace(matches[8]))

	if err != nil {
		log.Printf("Could not parse classNumber to int: %s\n", err.Error())
		return nil
	}

	teacherName := strings.TrimSpace(matches[1])
	lessionTitle := strings.TrimSpace(matches[3])
	lectureTitle := strings.TrimSpace(matches[5])
	subjectName := strings.TrimSpace(matches[6])
	subjectTypeParsed := strings.TrimSpace(matches[7])
	var subjectType db.SubjectType

	if subjectTypeParsed == "OLD" {
		subjectType = db.OLD_SYLLABUS
	} else {
		subjectType = db.NEW_SYLLABUS
	}

	video := &Video{
		TeacherName:   teacherName,
		LessonNumber:  lessionNumber,
		LessionName:   lessionTitle,
		LectureNumber: uint(lectureNumber),
		LectureName:   lectureTitle,
		SubjectName:   subjectName,
		SubjectType:   subjectType,
		Class:         uint(classNumber),
		VideoURL:      path,
	}

	return video
}
