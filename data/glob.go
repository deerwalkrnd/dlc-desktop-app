package data

import (
	"io/fs"
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
			videos = append(videos, video)
		}
		return nil
	})
	// fmt.Println(videos)
	return nil
}
