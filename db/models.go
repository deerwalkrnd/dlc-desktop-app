package db

import "gorm.io/gorm"

type SubjectType int

const (
	OLD_SYLLABUS SubjectType = iota
	NEW_SYLLABUS
)

type Class struct {
	gorm.Model
	Number   uint
	Subjects []Subject `gorm:"foreignKey:ClassId"`
}

type Subject struct {
	gorm.Model
	Name     string
	Type     SubjectType
	ClassId  uint
	Class    Class     `gorm:"foreignKey:ClassId"`
	Lectures []Lecture `gorm:"foreignKey:SubjectId"`
}

type Lecture struct {
	gorm.Model
	Number    uint
	Name      string
	SubjectId uint
	Subject   Subject  `gorm:"foreignKey:SubjectId"`
	Lessons   []Lesson `gorm:"foreignKey:LectureId"`
}

type Lesson struct {
	gorm.Model
	Name      string
	Number    float64
	VideoUrl  string
	TeacherId uint
	LectureId uint
	Lecture   Lecture `gorm:"foreignKey:LectureId"`
	Teacher   Teacher `gorm:"foreignKey:TeacherId"`
}

type Teacher struct {
	gorm.Model
	Name     string
	Lessions []Lesson `gorm:"foreignKey:TeacherId"`
}

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&Class{})
	db.AutoMigrate(&Subject{})
	db.AutoMigrate(&Lecture{})
	db.AutoMigrate(&Lesson{})
	db.AutoMigrate(&Teacher{})
}
