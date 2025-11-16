package api

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const LAYOUT = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", getNextDateHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.HandleFunc("/api/task/done", taskDoneHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		log.Printf("[REQ] error: wrong task handler method")
		writeErrorJSON(w, http.StatusMethodNotAllowed, "wrong method error")
		return
	}
}

func getNow(input string) (time.Time, error) {
	if input == "" {
		return time.Now(), nil
	} else {
		now, err := time.Parse(LAYOUT, input)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid 'now' format")
		}
		return now, nil
	}
}

func getDstart(input string) (time.Time, error) {
	if input == "" {
		return time.Time{}, fmt.Errorf("empty 'dstart' string")
	}

	output, err := time.Parse(LAYOUT, input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid 'dstart' format")
	}

	return output, nil
}

func checkRepeat(input string) error {
	// repeat := "m 5 1,13"
	parameters := strings.Split(input, " ")
	if len(parameters) < 1 || len(parameters) > 3 {
		return fmt.Errorf("incorrect number of 'repeat' parameters")
	}

	// repeatExpr := regexp.MustCompile(`^[ydmw](?:\s+\d+(?:,\d+)*)(?:\s+\d+(?:,\d+)*)?$`) // TODO
	repeatExpr := regexp.MustCompile(`^[yd](?:\s+\d+(?:,\d+)*(\s+\d+(?:,\d+)*)?)?$`)

	if repeatExpr.MatchString(input) {
		return nil
	}

	return fmt.Errorf("incorrect values of 'repeat' parameters")
}
