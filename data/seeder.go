package data

import (
	"fmt"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
	"gorm.io/gorm"
)

// SeedVideos seeds all videos inside a single transaction.
// If you want each video to be independent, move the transaction inside the loop.
func SeedVideos(videos []*Video, DB *gorm.DB) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		for i, video := range videos {
			if err := seedVideo(tx, video); err != nil {
				return fmt.Errorf("seed video[%d]: %w", i, err)
			}
		}
		return nil
	})
}

func seedVideo(tx *gorm.DB, v *Video) error {
	class, err := getOrCreateClass(tx, v.Class)
	if err != nil {
		return err
	}

	teacher, err := getOrCreateTeacher(tx, v.TeacherName)
	if err != nil {
		return err
	}

	subject, err := getOrCreateSubject(tx, class.ID, v.SubjectType, v.SubjectName)
	if err != nil {
		return err
	}

	unit, err := getOrCreateUnit(tx, subject.ID, v.UnitNumber, v.UnitName)
	if err != nil {
		return err
	}

	_, err = getOrCreateChapter(
		tx,
		unit.ID,
		teacher.ID,
		v.ChapterNumber,
		v.ChapterName,
		v.VideoURL,
	)
	if err != nil {
		return err
	}

	return nil
}

// ---------- helpers ----------

func getOrCreateClass(tx *gorm.DB, number uint) (*db.Class, error) {
	cl := db.Class{Number: number}
	if err := tx.Where("number = ?", number).FirstOrCreate(&cl).Error; err != nil {
		return nil, err
	}
	return &cl, nil
}

func getOrCreateTeacher(tx *gorm.DB, name string) (*db.Teacher, error) {
	t := db.Teacher{Name: name}
	if err := tx.Where("name = ?", name).FirstOrCreate(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func getOrCreateSubject(tx *gorm.DB, classID uint, subjectType db.SubjectType, name string) (*db.Subject, error) {
	s := db.Subject{}
	if err := tx.
		Where("name = ? AND type = ? AND class_id = ?", name, subjectType, classID).
		First(&s).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			s = db.Subject{
				Name:    name,
				Type:    subjectType,
				ClassId: classID,
			}
			if err := tx.Create(&s).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &s, nil
}

func getOrCreateUnit(tx *gorm.DB, subjectID, number uint, name string) (*db.Unit, error) {
	u := db.Unit{}
	if err := tx.
		Where("number = ? AND name = ? AND subject_id = ?", number, name, subjectID).
		First(&u).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			u = db.Unit{
				Number:    number,
				Name:      name,
				SubjectId: subjectID,
			}
			if err := tx.Create(&u).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &u, nil
}

func getOrCreateChapter(tx *gorm.DB, unitID, teacherID, number uint, name, videoURL string) (*db.Chapter, error) {
	c := db.Chapter{}
	if err := tx.
		Where("number = ? AND name = ? AND unit_id = ? AND teacher_id = ?", number, name, unitID, teacherID).
		First(&c).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			c = db.Chapter{
				Name:      name,
				Number:    number,
				VideoUrl:  videoURL,
				TeacherId: teacherID,
				UnitId:    unitID,
			}
			if err := tx.Create(&c).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &c, nil
}
