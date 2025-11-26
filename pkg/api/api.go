package api

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/igromanas/go-final/pkg/auth"
)

const LAYOUT = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", getNextDateHandler)
	// mux.HandleFunc("/api/task", taskHandler)
	// mux.HandleFunc("/api/tasks", tasksHandler)
	// mux.HandleFunc("/api/task/done", taskDoneHandler)
	mux.HandleFunc("/api/task", auth.Auth(taskHandler))
	mux.HandleFunc("/api/tasks", auth.Auth(tasksHandler))
	mux.HandleFunc("/api/task/done", auth.Auth(taskDoneHandler))
	mux.HandleFunc("/api/signin", signinHandler)
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

	var repeatExpr *regexp.Regexp
	switch parameters[0] {
	case "d":
		repeatExpr = regexp.MustCompile(`^d\s+0*(?:[1-9][0-9]{0,2}|400)$`)
	case "w":
		repeatExpr = regexp.MustCompile(`^w\s+0*[1-7](?:,0*[1-7])*$`)
	case "m":
		repeatExpr = regexp.MustCompile(`^m\s+(?:0*(?:[1-9]|[12][0-9]|3[01])|-1|-2)(?:,(?:0*(?:[1-9]|[12][0-9]|3[01])|-1|-2))*(?:\s+0*(?:[1-9]|1[0-2])(?:,0*(?:[1-9]|1[0-2]))*)?$`)
	case "y":
		return nil
	default:
		return fmt.Errorf("incorrect values of 'repeat' parameters in '%s", input)
	}
	// repeatExpr := regexp.MustCompile(`^[ydmw](?:\s+\d+(?:,\d+)*(\s+\d+(?:,\d+)*)?)?$`) // TODO
	// repeatExpr := regexp.MustCompile(`^[ydmw](?:\s+(?:-1|-2|[1-9]|[12]\d|3[01])(?:,(?:-1|-2|[1-9]|[12]\d|3[01]))*)?(?:\s+(?:[1-9]|1[0-2])(?:,(?:[1-9]|1[0-2]))*)?$`)
	// repeatExpr := regexp.MustCompile(`^[yd](?:\s+\d+(?:,\d+)*(\s+\d+(?:,\d+)*)?)?$`) // old

	if repeatExpr.MatchString(input) {
		return nil
	}

	return fmt.Errorf("incorrect values of 'repeat' parameters in '%s'", input)
}
