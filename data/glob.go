package data

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
)

func Initialize(dataPath string) error {
	videos := []*Video{}
	videoPattern := regexp.MustCompile(".mp4$")

	filepath.WalkDir(dataPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && videoPattern.MatchString(path) {
			video := ParseVideo(path)
			if video == nil {
				log.Printf("failed to parse: %s\n", path)
				return nil
			}
			videos = append(videos, video)
		}
		return nil
	})

	for _, video := range videos {
		fmt.Println(video)
	}

	return nil
}
