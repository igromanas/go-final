package api

import (
	"log"
	"net/http"

	"github.com/igromanas/go-final/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("[REQ] error: empty id")
		writeErrorJSON(w, http.StatusBadRequest, "empty id")
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		log.Printf("[DB] error getting task by id: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "db error")
		return
	}

	if task == nil {
		log.Printf("[REQ] task not found")
		writeErrorJSON(w, http.StatusNotFound, "no tasks")
		return
	}

	writeJSON(w, task)
}
