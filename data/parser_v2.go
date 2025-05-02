package data

import (
	"log"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
)

func ParseVideoV2(path string) *Video {
	basePath := filepath.Base(path)
	videoPath := basePath[0 : len(basePath)-4]
	items := strings.Split(videoPath, "-")

	lessionNumber, err := strconv.ParseFloat(strings.TrimSpace(items[1]), 64)
	if err != nil {
		log.Printf("Could not parse lessionNumber to float64: %s\n", err.Error())
		return nil
	}

	classNumber, err := strconv.Atoi(strings.TrimSpace(items[6]))
	if err != nil {
		log.Printf("Could not parse classNumber to int: %s\n", err.Error())
		return nil
	}

	teacherName := strings.TrimSpace(items[0])
	lessionTitle := strings.TrimSpace(items[2])
	lectureTitle := strings.TrimSpace(items[3])
	subjectName := strings.TrimSpace(items[4])
	subjectTypeParsed := strings.TrimSpace(items[5])
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
