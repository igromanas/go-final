package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/igromanas/go-final/pkg/auth"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[REQ] error: wrong method")
		writeErrorJSON(w, http.StatusMethodNotAllowed, "wrong method")
		return
	}

	var user LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("[REQ] error processing json: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, "json processing")
		return
	}

	token, err := auth.SignIn(user.Password)
	if err != nil {
		log.Printf("[REQ] incorrect password")
		writeErrorJSON(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	var resp LoginResponse
	resp.Token = token
	writeJSON(w, resp)
}
