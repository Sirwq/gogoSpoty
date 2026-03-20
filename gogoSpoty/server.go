package main

import (
	"gogoSpoty/spoty"
	"net/http"
)

func newServer(t *spoty.Track) *http.ServeMux {
	mux := http.NewServeMux()
	setupRoutes(mux, t)
	return mux
}

func setupRoutes(mux *http.ServeMux, t *spoty.Track) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/widget", spoty.WidgetHandler("static/widget.html"))
	mux.HandleFunc("/api/current", spoty.TrackHandler(t))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
