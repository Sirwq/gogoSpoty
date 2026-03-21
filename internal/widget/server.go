package widget

import (
	"net/http"
)

func NewServer(t *Track) *http.ServeMux {
	mux := http.NewServeMux()
	setupRoutes(mux, t)
	return mux
}

func setupRoutes(mux *http.ServeMux, t *Track) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/widget", WidgetHandler("static/widget.html"))
	mux.HandleFunc("/api/current", TrackHandler(t))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
