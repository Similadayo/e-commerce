package main

import (
	"log"
	"net/http"

	"github.com/Similadayo/backend/db"
	"github.com/Similadayo/backend/router"
)

func main() {
	//connect to database
	db.ConnectToDb()
	defer db.DB.Close()

	//
	r := router.Router()
	log.Fatal(http.ListenAndServe(":8080", r))
}
