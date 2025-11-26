package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/igromanas/go-final/pkg/db"
	"github.com/igromanas/go-final/pkg/model"
	"github.com/igromanas/go-final/service"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task model.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("[REQ] error processing json: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "json processing")
		return
	}

	if task.ID == "" {
		log.Printf("[REQ] error: empty id")
		writeErrorJSON(w, http.StatusBadRequest, "empty id")
		return
	}

	oldTask, err := db.GetTaskByID(task.ID)
	if err != nil {
		log.Printf("[DB] error getting task by id: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "db error")
		return
	}

	if oldTask == nil {
		log.Printf("[REQ] task '%s' not found", task.ID)
		writeErrorJSON(w, http.StatusNotFound, "no tasks")
		return
	}

	if task.Title == "" {
		log.Printf("[REQ] error: empty title")
		writeErrorJSON(w, http.StatusBadRequest, "empty title")
		return
	}

	if service.ModifyDate(&task) != nil {
		log.Printf("[REQ] error of incorrect date: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "incorrect date")
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		log.Printf("[DB] error updating task '%s': %v", task.ID, err)
		writeErrorJSON(w, http.StatusBadRequest, "db error")
		return
	}

	writeJSON(w, struct{}{})
}
