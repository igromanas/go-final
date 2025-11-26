package server

import (
	"net/http"

	"github.com/igromanas/go-final/pkg/api"
)

func Run(port string) error {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	api.Init(mux)

	return http.ListenAndServe(port, mux)
}
