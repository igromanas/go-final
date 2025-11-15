package server

import (
	"log"
	"net/http"
	"os"

	"github.com/igromanas/go-final/pkg/api"
)

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":7540"
	}
	return port
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./web/index.html")
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func Run() {
	mux := http.NewServeMux()

	// mux.Handle("/js/", http.FileServer(http.Dir("./web/js")))
	// mux.Handle("/css/", http.FileServer(http.Dir("./web/css")))
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	// mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	// mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	// mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	// mux.HandleFunc("/", handlerFunc)

	api.Init(mux)

	port := getPort()
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
