package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/deerwalkrnd/dlc-desktop-app/api"
	"github.com/deerwalkrnd/dlc-desktop-app/data"
	db "github.com/deerwalkrnd/dlc-desktop-app/db"
)

var Logger = log.Default()

func init() {
	Logger.Println("Started DLC Desktop Application")

	if _, err := os.Stat(db.DATABASE_NAME); err != nil {
		// no db file detected, create - migrate - populate
		DB, err := db.GetDB()

		if err != nil {
			Logger.Fatalf("error getting database, %s\n", err.Error())
		}

		db.MigrateModels(DB)

		Logger.Println("Database Migration Finished")

		dataPath, _ := filepath.Abs("DLC")

		err = data.Initialize(dataPath, DB)

		if err != nil {
			Logger.Fatalf("failed to initialize the database and seed data: %s\n", err.Error())
		}

	}
}
func main() {
	outputPath := "./web/build"

	mainMux := http.NewServeMux()
	apiMux := api.GetApiMux()

	fs := http.FileServer(http.Dir(outputPath))

	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(outputPath, r.URL.Path)
		_, err := os.Stat(path)

		if os.IsNotExist(err) && filepath.Ext(r.URL.Path) == "" {
			http.ServeFile(w, r, filepath.Join(outputPath, "index.html"))
			return
		}

		fs.ServeHTTP(w, r)
	})

	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mainMux)
	if err != nil {
		log.Fatal(err)
	}
}
