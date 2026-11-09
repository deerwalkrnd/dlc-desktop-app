// Command dlc-check validates a DLC video folder before it is shipped. It
// lists every file that the app would skip or that looks like a mistake,
// using exactly the same parser as the app.
//
// Usage:
//
//	dlc-check [-v] [folder]
//
// The folder defaults to ./DLC. The exit status is 1 if any file would be
// skipped, so it can be used in scripts.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/deerwalkrnd/dlc-desktop-app/internal/library"
)

func main() {
	verbose := flag.Bool("v", false, "also list every lesson that was found")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-v] [folder]\n\nChecks that every .mp4 in the folder (default ./DLC) is named correctly.\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	root := "DLC"
	if flag.NArg() > 0 {
		root = strings.Join(flag.Args(), " ") // tolerate unquoted paths with spaces
	}

	x, err := library.Build(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	if !x.Exists {
		fmt.Fprintf(os.Stderr, "error: folder %q does not exist\n", root)
		os.Exit(2)
	}

	if *verbose {
		for _, c := range x.Classes {
			fmt.Printf("Class %d\n", c.Number)
			for _, s := range c.Subjects {
				fmt.Printf("  %s (%s syllabus)\n", s.Name, s.Syllabus)
				for _, ch := range s.Chapters {
					fmt.Printf("    Chapter %d: %s\n", ch.Number, ch.Title)
					for _, l := range ch.Lessons {
						fmt.Printf("      %s  %s  [%s]\n", l.Number, l.Title, l.Teacher)
					}
				}
			}
		}
		fmt.Println()
	}

	if len(x.Warnings) > 0 {
		fmt.Println("Warnings (these files are included, but check them):")
		for _, w := range x.Warnings {
			fmt.Printf("  %s\n      %s\n", w.Path, w.Message)
		}
		fmt.Println()
	}
	if len(x.Skipped) > 0 {
		fmt.Println("Skipped (these files will NOT appear in the app):")
		for _, s := range x.Skipped {
			fmt.Printf("  %s\n      %s\n", s.Path, s.Message)
		}
		fmt.Println()
		fmt.Println(`Expected: TEACHER - LESSON NO - CHAPTER TITLE - LESSON TITLE - SUBJECT - OLD|NEW - CLASS.mp4`)
		fmt.Println()
	}

	fmt.Printf("%d lessons in %d classes, %d warnings, %d skipped\n",
		x.LessonCount(), len(x.Classes), len(x.Warnings), len(x.Skipped))
	if len(x.Skipped) > 0 {
		os.Exit(1)
	}
}
