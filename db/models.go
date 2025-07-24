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
	Name    string
	Type    SubjectType
	ClassId uint
	Class   Class  `gorm:"foreignKey:ClassId"`
	Units   []Unit `gorm:"foreignKey:SubjectId"`
}

type Unit struct {
	gorm.Model
	Number    uint
	Name      string
	SubjectId uint
	Subject   Subject   `gorm:"foreignKey:SubjectId"`
	Chapters  []Chapter `gorm:"foreignKey:UnitId"`
}

type Chapter struct {
	gorm.Model
	Name      string
	Number    uint
	VideoUrl  string
	TeacherId uint
	UnitId    uint
	Unit      Unit    `gorm:"foreignKey:UnitId"`
	Teacher   Teacher `gorm:"foreignKey:TeacherId"`
}

type Teacher struct {
	gorm.Model
	Name     string
	Chapters []Chapter `gorm:"foreignKey:TeacherId"`
}

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&Class{})
	db.AutoMigrate(&Subject{})
	db.AutoMigrate(&Unit{})
	db.AutoMigrate(&Chapter{})
	db.AutoMigrate(&Teacher{})
}
