package data

import (
	"log"
	"strconv"
	"strings"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
	"github.com/deerwalkrnd/dlc-desktop-app/lexer"
)

func ParseVideoV3(path string) *Video {
	lexer_ins := lexer.NewVideoLexer(path)

	if err := lexer_ins.Tokenize(); err != nil {
		log.Println("failed to tokenize filename " + path + err.Error())
		return nil
	}

	if err := lexer_ins.Validate(); err != nil {
		log.Println("Invalid Filename Format " + path + err.Error())
		return nil
	}

	teacherToken, _ := lexer_ins.GetTokenByType(lexer.TEACHER_NAME)
	teacherName := strings.TrimSpace(teacherToken.Value)

	lessonToken, _ := lexer_ins.GetTokenByType(lexer.UNIT_NUMBER)
	lessonNumberStr := strings.TrimSpace(lessonToken.Value)

	lessonParts := strings.Split(lessonNumberStr, ".")

	if len(lessonParts) != 2 {
		log.Println("invalid lesson parts " + lessonNumberStr)
		return nil
	}
	unitNumber, err := strconv.Atoi(lessonParts[0])

	if err != nil {
		log.Printf("Could not parse unit snumber: %s", err.Error())
		return nil
	}

	chapterNumber, err := strconv.ParseFloat(lessonParts[1], 64)

	if err != nil {
		log.Printf("Could not parse sub-lesson/chapter number: %s", err.Error())
		return nil
	}

	fullLessonNumber, err := strconv.ParseFloat(lessonNumberStr, 64)
	if err != nil {
		log.Printf("Could not parse full lesson number: %s", err.Error())
		return nil
	}

	unitToken, _ := lexer_ins.GetTokenByType(lexer.UNIT_NAME)
	unitTitle := strings.TrimSpace(unitToken.Value)

	chapterToken, _ := lexer_ins.GetTokenByType(lexer.CHAPTER_NAME)
	chapterTitle := strings.TrimSpace(chapterToken.Value)

	subjectToken, _ := lexer_ins.GetTokenByType(lexer.SUBJECT_NAME)
	subjectName := strings.TrimSpace(subjectToken.Value)

	// Extract and parse subject type
	subjectTypeToken, _ := lexer_ins.GetTokenByType(lexer.SUBJECT_TYPE)
	subjectTypeParsed := strings.TrimSpace(subjectTypeToken.Value)

	var subjectType db.SubjectType

	if strings.ToUpper(subjectTypeParsed) == "OLD" {
		subjectType = db.OLD_SYLLABUS
	} else {
		subjectType = db.NEW_SYLLABUS
	}

	// Extract and parse class number
	classToken, _ := lexer_ins.GetTokenByType(lexer.CLASS_NUMBER)

	classNumber, err := strconv.Atoi(strings.TrimSpace(classToken.Value))
	if err != nil {
		log.Printf("Could not parse class number: %s", err.Error())
		return nil
	}

	return &Video{
		Class:            uint(classNumber),
		SubjectType:      subjectType,
		SubjectName:      subjectName,
		TeacherName:      teacherName,
		UnitNumber:       uint(unitNumber),
		UnitName:         unitTitle,
		ChapterNumber:    uint(chapterNumber),
		ChapterName:      chapterTitle,
		FullLessonNumber: fullLessonNumber,
		VideoURL:         path,
	}

	return nil
}
