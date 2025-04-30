package data

import (
	"io/fs"
	"log"
	"path/filepath"
	"regexp"

	"gorm.io/gorm"
)

var VideoPattern = regexp.MustCompile(".mp4$")

func Initialize(dataPath string, db *gorm.DB) error {
	videos := []*Video{}

	filepath.WalkDir(dataPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && VideoPattern.MatchString(path) {
			video := ParseVideo(path)
			if video == nil {
				log.Printf("failed to parse: %s\n", path)
				log.Printf("skipping!")
				return nil
			}
			videos = append(videos, video)
		}
		return nil
	})

	err := SeedVideos(videos, db)

	if err != nil {
		return err
	}

	return nil
}
