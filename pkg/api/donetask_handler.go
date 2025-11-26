package api

import (
	"log"
	"net/http"
	"time"

	"github.com/igromanas/go-final/pkg/db"
	"github.com/igromanas/go-final/service"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[REQ] error: wrong method")
		writeErrorJSON(w, http.StatusMethodNotAllowed, "wrong method error")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("[REQ] error: empty ID")
		writeErrorJSON(w, http.StatusBadRequest, "empty id error")
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

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			log.Printf("[DB] error deleting task '%s': %v", id, err)
			writeErrorJSON(w, http.StatusInternalServerError, "db error")
			return
		}
		writeJSON(w, struct{}{})
		return
	}

	dstart, err := getDstart(task.Date)
	if err != nil {
		log.Printf("[ERR] error converting task.date: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "date conversion error")
		return
	}

	var nd time.Time
	if time.Now().Format(LAYOUT) <= task.Date {
		nd, err = service.NextDate(time.Now(), dstart, task.Repeat)
		if err != nil {
			log.Printf("[ERR] error processing next day: %v", err)
			writeErrorJSON(w, http.StatusInternalServerError, "processing next day error")
			return
		}
	}

	err = db.UpdateDate(task.ID, nd.Format(LAYOUT))
	if err != nil {
		log.Printf("[DB] error updating task: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, struct{}{})
}
