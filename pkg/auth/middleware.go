package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeErrorJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	respError := map[string]any{
		"error": msg,
	}

	err := json.NewEncoder(w).Encode(respError)
	if err != nil {
		log.Printf("error processing json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			log.Printf("[REQ] incorrect token")
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required")
			// http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		jwt := cookie.Value

		if err := ValidateJWT(jwt); err != nil {
			log.Printf("[REQ] incorrect password")
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required")
			// http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
