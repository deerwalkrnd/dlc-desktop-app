package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const DATABASE_NAME = "dlc.sqlite"

func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DATABASE_NAME), &gorm.Config{})
	return db, err
}
