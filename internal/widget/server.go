package widget

import (
	"embed"
	"net/http"
)

func NewServerMux(t *Track) *http.ServeMux {
	mux := http.NewServeMux()
	setupRoutes(mux, t)
	return mux
}

//go:embed static/*
var staticFiles embed.FS

func setupRoutes(mux *http.ServeMux, t *Track) {
	mux.Handle("/static/", http.FileServerFS(staticFiles))
	mux.HandleFunc("/api/current", TrackHandler(t))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/widget", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/widget.html", http.StatusFound)
	})
}
