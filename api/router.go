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

	ApiMux.HandleFunc("GET /teachers", apiHandler.GetTeachers)
	ApiMux.HandleFunc("GET /classes", apiHandler.GetClasses)
	ApiMux.HandleFunc("GET /class", apiHandler.GetLecturesByClass)

	return ApiMux
}
