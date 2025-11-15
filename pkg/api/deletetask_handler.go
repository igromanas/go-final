package api

import (
	"log"
	"net/http"

	"github.com/igromanas/go-final/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		log.Printf("[REQ] error: wrong method")
		writeErrorJSON(w, http.StatusMethodNotAllowed, "wrong method")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("[REQ] error: empty ID")
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
		log.Printf("[REQ] task '%s' not found", id)
		writeErrorJSON(w, http.StatusNotFound, "no tasks")
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		log.Printf("[DB] error deleting task '%s': %v", id, err)
		writeErrorJSON(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, struct{}{})
}
