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

	UnitNumber uint
	UnitName   string

	ChapterNumber uint

	ChapterName      string
	FullLessonNumber float64

	VideoURL string
}

func NewVideo(
	TeacherName, ChapterName, UnitName, SubjectName string,
	ChapterNumber uint, UnitNumber, Class uint, FullLessonNumber float64,
	SubjectType db.SubjectType, VideoURL string) *Video {

	return &Video{
		TeacherName:      TeacherName,
		UnitNumber:       UnitNumber,
		UnitName:         UnitName,
		ChapterNumber:    ChapterNumber,
		ChapterName:      ChapterName,
		FullLessonNumber: FullLessonNumber,
		SubjectName:      SubjectName,
		SubjectType:      SubjectType,
		VideoURL:         VideoURL,
		Class:            Class,
	}
}
