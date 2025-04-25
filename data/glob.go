package data

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
)

func Initialize(dataPath string) error {
	videos := []string{}
	videoPattern := regexp.MustCompile(".mp4$")

	filepath.WalkDir(dataPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && videoPattern.MatchString(path) {
			fmt.Println(path)
			videos = append(videos, path)
		}
		return nil
	})
	fmt.Println(videos)
	return nil
}
