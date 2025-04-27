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

	fmt.Println("got class : ", *class)
	fmt.Println("got teacher : ", *teacher)

	return nil
}

func getClass(classNumber int, DB *gorm.DB) *db.Class {
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
