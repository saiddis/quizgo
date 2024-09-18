package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"gihub.com/saiddis/quizgo/internal/install/database"
)

type apiConfig struct {
	DB *database.Queries
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with %d error: %s", code, msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
