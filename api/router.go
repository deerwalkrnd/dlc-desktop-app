package api

import (
	"log"
	"net/http"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
)

func GetApiMux() *http.ServeMux {

	var ApiMux = http.NewServeMux()
	db, err := db.GetDB()

	if err != nil {
		log.Fatalln("error: ", err.Error())
	}

	apiHandler := NewApiHandler(db)

	ApiMux.HandleFunc("/teachers", apiHandler.GetTeachers)
	ApiMux.HandleFunc("/classes", apiHandler.GetClasses)

	return ApiMux
}
