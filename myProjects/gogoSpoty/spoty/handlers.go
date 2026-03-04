package spoty

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func TrackHandler(t *Track) http.HandlerFunc {
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

func WidgetHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

func CallbackHandler(state string, auth *spotifyauth.Authenticator, ch chan *oauth2.Token) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Auth failed", http.StatusForbidden)
			return
		}
		fmt.Fprintf(w, "Login success!")
		ch <- token
	}
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
