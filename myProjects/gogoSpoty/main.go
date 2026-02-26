package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

func main() {
	godotenv.Load(".env")

	port := ":5111"
	redirUrl := "http://127.0.0.1:5111/callback"

	var t Track
	var waitTime time.Duration = 5
	state, auth, ch := OAuthFlow(redirUrl)

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/widget", widgetHandler())
	mux.HandleFunc("/api/current", trackHandler(&t))
	mux.HandleFunc("/callback", callbackHandler(state, auth, ch))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	go http.ListenAndServe(port, mux)
	ctx := context.Background()
	token, err := loadToken()
	if err != nil {
		token = <-ch
		saveToken(token)
	}

	client := spotify.New(auth.Client(ctx, token))

	fmt.Println("Server is running...")

	go func() {
		for {
			playing, err := client.PlayerCurrentlyPlaying(ctx)

			if err != nil {
				log.Printf("Error fetching track %v", err)
				time.Sleep(waitTime * time.Second)
				continue
			}

			if playing != nil {
				updateTrack(&t, playing)
			}

			time.Sleep(waitTime * time.Second)
		}
	}()

	select {}

}
