package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/igromanas/go-final/service"
)

func getNextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("[REQ] error: wrong method")
		writeErrorJSON(w, http.StatusMethodNotAllowed, "wrong method error")
		return
	}

	now, err := getNow(r.URL.Query().Get("now"))
	if err != nil {
		log.Printf("[REQ] error getting NOW: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "'now' error")
		return
	}

	dstart, err := getDstart(r.URL.Query().Get("date"))
	if err != nil {
		log.Printf("[REQ] error converting date: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "'dstart' error")
		return
	}

	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		writeJSON(w, struct{}{})
		return
	}

	err = checkRepeat(repeat)
	if err != nil {
		log.Printf("[REQ] error checking repeat (%s): %v", repeat, err)
		writeErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("'repeat' (%s) error", repeat))
		return
	}

	nd, err := service.NextDate(now, dstart, repeat)
	if err != nil {
		log.Printf("[REQ] error getting next day: %v", err)
		writeErrorJSON(w, http.StatusUnprocessableEntity, "'repeat' rule error")
		return
	}

	log.Printf("NextDay result: %s\n", nd.Format(LAYOUT))

	_, err = w.Write([]byte(nd.Format(LAYOUT)))
	if err != nil {
		log.Printf("[ERR] error writing next day: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "writing response")
		return
	}
}
