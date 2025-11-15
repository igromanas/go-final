package server

import (
	"log"
	"net/http"

	"github.com/igromanas/go-final/pkg/api"
)

func Run(port string) {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	api.Init(mux)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
