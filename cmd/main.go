package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	outputPath := "./web/build"

	fs := http.FileServer(http.Dir(outputPath))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(outputPath, r.URL.Path)
		_, err := os.Stat(path)

		if os.IsNotExist(err) && filepath.Ext(r.URL.Path) == "" {
			http.ServeFile(w, r, filepath.Join(outputPath, "index.html"))
			return
		}

		fs.ServeHTTP(w, r)
	})

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
