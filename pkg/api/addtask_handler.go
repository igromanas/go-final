package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/igromanas/go-final/pkg/db"
	"github.com/igromanas/go-final/pkg/model"
	"github.com/igromanas/go-final/service"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("[REQ] error processing json: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "json processing")
		return
	}

	if task.Title == "" {
		log.Printf("[REQ] error: empty title")
		writeErrorJSON(w, http.StatusBadRequest, "empty title")
		return
	}

	if err := service.ModifyDate(&task); err != nil {
		log.Printf("[REQ] error of incorrect date: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "incorrect date")
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("[DB] error adding task: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "db error")
		return
	}

	resp := map[string]any{
		"id": id,
	}
	writeJSON(w, resp)
}
