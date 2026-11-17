package server

import "net/http"

func SetHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /", defaultPath)
}

func defaultPath(w http.ResponseWriter, r *http.Request) {
}
