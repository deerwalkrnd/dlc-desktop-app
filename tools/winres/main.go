// Command winres writes a .syso file with the Windows icon, version
// information and manifest for an executable. The Go linker picks it up
// automatically when building for Windows. `make release` runs it.
//
//	go run ./tools/winres -version v2.0.0 -description "DLC Desktop App" \
//	    -original dlc-desktop-app.exe -out cmd/dlc-desktop-app/rsrc_windows_amd64.syso
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deerwalkrnd/dlc-desktop-app/internal/winres"
)

func main() {
	var (
		out         = flag.String("out", "", "output .syso file (required)")
		arch        = flag.String("arch", "amd64", "target GOARCH: amd64 or arm64")
		version     = flag.String("version", "dev", "version string, e.g. v2.0.0")
		description = flag.String("description", "DLC Desktop App", "file description shown by Windows")
		original    = flag.String("original", "", "original file name, e.g. dlc-desktop-app.exe")
		icons       = flag.String("icons", "assets/windows", "folder with icon-*.png files")
	)
	flag.Parse()
	if *out == "" || flag.NArg() > 0 {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*out, *arch, *icons, winres.Resources{
		Version:         *version,
		ProductName:     "DLC Desktop App",
		FileDescription: *description,
		CompanyName:     "Deerwalk Learning Center",
		Copyright:       fmt.Sprintf("Copyright © %d Deerwalk Learning Center. MIT License.", time.Now().Year()),
		OriginalName:    *original,
	}); err != nil {
		fmt.Fprintln(os.Stderr, "winres:", err)
		os.Exit(1)
	}
}

func run(out, arch, iconDir string, r winres.Resources) error {
	files, err := filepath.Glob(filepath.Join(iconDir, "icon-*.png"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no icon-*.png files in %s", iconDir)
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		r.Icons = append(r.Icons, data)
	}
	syso, err := winres.Build(r, arch)
	if err != nil {
		return err
	}
	return os.WriteFile(out, syso, 0o644)
}
