package data

import "github.com/deerwalkrnd/dlc-desktop-app/db"

// TEACHER NAME - LESSION NO | LESSION TITLE - LECTURE NO | LECTURE TITLE - SUBJECT NAME | SUBJECT TYPE - CLASS .mp4

type Video struct {
	TeacherName  string
	LessonNumber float64

	LessionTitle  string
	LectureNumber int

	LectureTitle string
	SubjectName  string

	SubjectType db.SubjectType
}

func NewVideo(TeacherName, LessonTitle, LectureTitle, SubjectName string, LessionNumber float64, LectureNumber int, SubjectType db.SubjectType) *Video {
	return &Video{
		TeacherName:   TeacherName,
		LessonNumber:  LessionNumber,
		LessionTitle:  LectureTitle,
		LectureNumber: LectureNumber,
		LectureTitle:  LectureTitle,
		SubjectName:   SubjectName,
		SubjectType:   SubjectType,
	}
}
