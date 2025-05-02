package data

import (
	"fmt"
	"log"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
)

func ParseVideo(path string) *Video {
	basePath := filepath.Base(path)
	fmt.Println("parsing: ", basePath)

	re := regexp.MustCompile(MATCH_VIDEO_PATTERN)
	matches := re.FindStringSubmatch(basePath)

	if matches == nil || len(matches) != 8 { // 7 capture groups + 1 for the whole match = 8
		log.Printf("Failed to parse, not in correct format: %s\n", path)
		return nil
	}

	lessionNumber, err := strconv.ParseFloat(strings.TrimSpace(matches[2]), 64)
	if err != nil {
		log.Printf("Could not parse lessionNumber to float64: %s\n", err.Error())
		return nil
	}

	classNumber, err := strconv.Atoi(strings.TrimSpace(matches[7]))
	if err != nil {
		log.Printf("Could not parse classNumber to int: %s\n", err.Error())
		return nil
	}

	teacherName := strings.TrimSpace(matches[1])
	lessionTitle := strings.TrimSpace(matches[3])
	lectureTitle := strings.TrimSpace(matches[4])
	subjectName := strings.TrimSpace(matches[5])
	subjectTypeParsed := strings.TrimSpace(matches[6])
	var subjectType db.SubjectType

	if subjectTypeParsed == "OLD" {
		subjectType = db.OLD_SYLLABUS
	} else {
		subjectType = db.NEW_SYLLABUS
	}

	lectureNumber := math.Floor(lessionNumber)

	video := &Video{
		TeacherName:   teacherName,
		LessonNumber:  lessionNumber,
		LessionName:   lessionTitle,
		LectureName:   lectureTitle,
		SubjectName:   subjectName,
		LectureNumber: uint(lectureNumber),
		SubjectType:   subjectType,
		Class:         uint(classNumber),
		VideoURL:      basePath,
	}

	return video
}
