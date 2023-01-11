package main

import (
	"log"
	"net/http"

	"github.com/Similadayo/backend/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := router.Router()
	log.Fatal(http.ListenAndServe(":8080", r))
}
