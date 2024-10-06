package respose

import (
	"log"
	"net/http"
)

func WithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with %d error: %s", code, msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	WithJSON(w, code, errResponse{
		Error: msg,
	})
}
