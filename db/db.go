package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DATABASE_NAME = "dlc.sqlite"

func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DATABASE_NAME), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	return db, err
}
