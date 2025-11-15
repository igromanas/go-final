package api

import (
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
		log.Printf("[REQ] error getting NOW: {%v}", err)
		writeErrorJSON(w, http.StatusBadRequest, "'now' error")
		return
	}

	dstart, err := getDstart(r.URL.Query().Get("date"))
	if err != nil {
		log.Printf("[REQ] error converting date: {%v}", err)
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
		log.Printf("[REQ] error checking repeat: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "'repeat' error")
		return
	}

	nd, err := service.NextDate(now, dstart, repeat)
	if err != nil {
		log.Printf("[REQ] error getting next day: {%v}", err)
		writeErrorJSON(w, http.StatusUnprocessableEntity, "'repeat' rule error")
		return
	}

	log.Printf("result: %s\n", nd.Format(LAYOUT))
	// resp := map[string]any{
	// 	"next_date": nd.Format(LAYOUT),
	// }

	// writeJSON(w, resp)
	writeJSON(w, nd.Format(LAYOUT))
}
