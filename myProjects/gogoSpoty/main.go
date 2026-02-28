package main

import (
	"context"
	"fmt"
	"gogoSpoty/spoty"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

/* implement .ENV loading trough UI */

func main() {
	godotenv.Load(".env")

	port := ":5111"
	redirUrl := "http://127.0.0.1:5111/callback"

	var t spoty.Track
	var waitTime time.Duration = 5
	state, auth, ch := spoty.OAuthFlow(redirUrl)
	mux := http.NewServeMux()
	setupRoutes(mux, &t, state, auth, ch)

	go http.ListenAndServe(port, mux)
	ctx := context.Background()
	token, err := spoty.LoadToken()
	if err != nil {
		url := auth.AuthURL(state)
		fmt.Println("Open this url: ", url)
		token = <-ch
		spoty.SaveToken(token)
	}

	client := spotify.New(auth.Client(ctx, token))

	fmt.Println(time.Now().Clock())
	fmt.Println("Server is running...")

	go func() {
		for {
			playing, err := client.PlayerCurrentlyPlaying(ctx)

			if err != nil {
				log.Printf("Error fetching track %v", err)
				time.Sleep(waitTime * time.Second)
				continue
			}

			if playing != nil && playing.Item != nil {
				spoty.UpdateTrack(&t, playing)
			}

			time.Sleep(waitTime * time.Second)
		}
	}()

	select {}

}

func setupRoutes(mux *http.ServeMux, t *spoty.Track, state string, auth *spotifyauth.Authenticator, ch chan *oauth2.Token) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/widget", spoty.WidgetHandler())
	mux.HandleFunc("/api/current", spoty.TrackHandler(t))
	mux.HandleFunc("/callback", spoty.CallbackHandler(state, auth, ch))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
