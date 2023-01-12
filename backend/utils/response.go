package utils

import (
	"encoding/json"
	"net/http"
)

// Respond is a helper function to format and return a response to the client
func Respond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
