package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error processing json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeErrorJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	respError := map[string]any{
		"error": msg,
	}

	err := json.NewEncoder(w).Encode(respError)
	if err != nil {
		log.Printf("error processing json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
