package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/deerwalkrnd/dlc-desktop-app/data"
)

func main() {
	var userInput string
	fmt.Print("Enter Directory Path: ")
	fmt.Scanln(&userInput)

	fmt.Println("Directory: " + userInput)
	var videos []*data.Video
	var unparsed []string

	err := filepath.WalkDir(userInput, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if !d.IsDir() && data.VideoPattern.MatchString(path) {
			// fmt.Println(path)
			video := data.ParseVideoV2(path)
			if video == nil {
				log.Printf("failed to parse: %s\n", path)
				unparsed = append(unparsed, path)
				log.Printf("")
				return nil
			}
			videos = append(videos, video)
		}
		return nil
	})

	if err != nil {
		log.Fatalln("Error: ", err.Error())
	}

	fmt.Println("Parsed Videos: ")

	for _, video := range videos {
		fmt.Println(video.VideoURL)
	}

	fmt.Println("Parsed Count: ", len(videos))
	if len(unparsed) > 0 {
		fmt.Println("Failed to Parse: ")
		for _, video := range unparsed {
			fmt.Println(filepath.Base(video))
		}
		fmt.Println("Failed Count: ", len(unparsed))
	}

}
