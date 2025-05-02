package data

import (
	"github.com/deerwalkrnd/dlc-desktop-app/db"
)

// TEACHER NAME - LESSION NO - LESSION TITLE - LECTURE NO - LECTURE TITLE - SUBJECT NAME - SUBJECT TYPE - CLASS .mp4
const MATCH_VIDEO_PATTERN = `^([^-]+) - ([^-]+) - ([^-]+) - ([^-]+) - ([^-]+) - ([^-]+) - (\d+) \.mp4$`

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
