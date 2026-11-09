// Package web embeds the built SvelteKit app so the desktop binary is a
// single self-contained file.
//
// The build/ directory is produced by `npm run build` in this folder (or
// `make build` at the repository root). A placeholder file keeps the embed
// pattern valid in a fresh checkout; the server then shows instructions
// instead of the app.
package web

import (
	"embed"
	"io/fs"
)

//go:embed all:build
var build embed.FS

// Dist returns the built web app, rooted at the directory holding index.html.
func Dist() fs.FS {
	sub, err := fs.Sub(build, "build")
	if err != nil {
		panic(err) // unreachable: "build" is a valid, embedded path
	}
	return sub
}
