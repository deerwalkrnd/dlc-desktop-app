package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/deerwalkrnd/dlc-desktop-app/api"
	"github.com/deerwalkrnd/dlc-desktop-app/data"
	db "github.com/deerwalkrnd/dlc-desktop-app/db"
	"github.com/gorilla/mux"
)

var Logger = log.Default()

const DATA_FOLDER string = "DLC"

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

		dataPath, _ := filepath.Abs(DATA_FOLDER)

		err = data.Initialize(dataPath, DB)

		if err != nil {
			Logger.Fatalf("failed to initialize the database and seed data: %s\n", err.Error())
		}

	}
}
func main() {
	outputPath := "./web/build"

	mainRouter := mux.NewRouter()

	mainRouter.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	mainRouter.PathPrefix("/videos/").Handler(http.StripPrefix("/videos/",
		http.FileServer(http.Dir("./"+DATA_FOLDER)))).Methods("GET")

	spaHandler := api.SpaHandler{
		StaticPath: outputPath,
		IndexPath:  "index.html",
	}

	apiHandler := api.GetApiRouter()

	apiHandler.SetupRoutes(mainRouter)
	mainRouter.PathPrefix("/").Handler(spaHandler)

	srv := &http.Server{
		Handler:      mainRouter,
		Addr:         "127.0.0.1:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Listening on :3000...")
	log.Fatal(srv.ListenAndServe())

}
