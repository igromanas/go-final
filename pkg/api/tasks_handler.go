package api

import (
	"log"
	"net/http"

	"github.com/igromanas/go-final/pkg/db"
	"github.com/igromanas/go-final/pkg/model"
)

const TaskLimit = 50

type TasksResp struct {
	Tasks []*model.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("[REQ] error: wrong method")
		writeErrorJSON(w, http.StatusMethodNotAllowed, "wrong method error")
		return
	}

	tasks, err := db.Tasks(TaskLimit)
	if err != nil {
		log.Printf("[DB] error getting tasks: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "db error")
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
