package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

}
