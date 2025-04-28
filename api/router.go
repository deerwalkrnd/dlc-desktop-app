package api

import (
	"log"

	"github.com/deerwalkrnd/dlc-desktop-app/db"
)

func GetApiRouter() *ApiHandler {

	db, err := db.GetDB()

	if err != nil {
		log.Fatalln("error: ", err.Error())
	}

	apiHandler := NewApiHandler(db)

	return apiHandler
}
