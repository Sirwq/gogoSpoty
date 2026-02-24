package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func trackHandler(t *Track) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.mx.Lock()
		defer t.mx.Unlock()
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(t)

		if err != nil {
			http.Error(w, "StatusUnprocessableEntity", http.StatusUnprocessableEntity)
			check(err, "marshaling json")
			return
		}
		w.Write(data)
	}
}

func widgetHandler(t *Track) http.HandlerFunc {
	tmpl, err := template.ParseFiles("templates/index.html")
	check(err, "parsing template")

	return func(w http.ResponseWriter, r *http.Request) {
		data := t
		err = tmpl.Execute(w, data)
		check(err, "executing template")
	}

}

func check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
