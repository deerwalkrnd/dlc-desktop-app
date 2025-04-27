package data

import (
	"fmt"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
	"gorm.io/gorm"
)

func SeedVideos(videos []*Video, DB *gorm.DB) error {

	for _, video := range videos {
		seedVideo(video, DB)
	}

	return nil
}

func seedVideo(video *Video, db *gorm.DB) error {
	class := getClass(video.Class, db)
	teacher := getTeacher(video.TeacherName, db)
	subject := getSubject(class.ID, video.SubjectType, video.SubjectName, db)

	lecture := getLecture(
		video.LectureNumber,
		video.LectureName,
		subject.ID,
		db,
	)

	fmt.Println("got class : ", *class)
	fmt.Println("got teacher : ", *teacher)
	fmt.Println("got subject: ", *subject)
	fmt.Println("got lecture: ", *lecture)
	return nil
}

func getClass(classNumber uint, DB *gorm.DB) *db.Class {

	var class db.Class
	result := DB.Where("number = ?", classNumber).First(&class)

	if result.Error == gorm.ErrRecordNotFound {
		class = db.Class{Number: classNumber}
		DB.Create(&class)

	}

	return &class
}

func getTeacher(teacherName string, DB *gorm.DB) *db.Teacher {

	var teacher db.Teacher
	result := DB.Where("name = ?", teacherName).First(&teacher)

	if result.Error == gorm.ErrRecordNotFound {

		teacher = db.Teacher{Name: teacherName}
		DB.Create(&teacher)

	}

	return &teacher
}

func getSubject(classId uint, subjectType db.SubjectType, subjectName string, DB *gorm.DB) *db.Subject {

	var subject db.Subject
	result := DB.Where("name = ? AND type = ?", subjectName, subjectType).First(&subject)

	if result.Error == gorm.ErrRecordNotFound {

		subject = db.Subject{
			Name:    subjectName,
			Type:    subjectType,
			ClassId: classId,
		}

		DB.Create(&subject)

	}

	return &subject
}

func getLecture(lectureNumber uint, lectureName string, subjectId uint, DB *gorm.DB) *db.Lecture {
	var lecture db.Lecture
	result := DB.Where("number = ? AND name = ? AND subject_id = ?", lectureNumber, lectureName, subjectId).First(&lecture)

	if result.Error == gorm.ErrRecordNotFound {
		lecture = db.Lecture{
			Number:    lectureNumber,
			Name:      lectureName,
			SubjectId: subjectId,
		}
		DB.Create(&lecture)

	}

	return &lecture
}
