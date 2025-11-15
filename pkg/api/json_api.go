package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := json.Marshal(data)
	if err != nil {
		log.Printf("error processing json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func writeErrorJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	respError := map[string]any{
		"error": msg,
	}
	resp, err := json.Marshal(respError)
	if err != nil {
		log.Printf("error processing json: %v", err)
		w.WriteHeader(code)
		return
	}
	w.Write(resp)
}
