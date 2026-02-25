package main

import (
	"context"
	"fmt"
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

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/widget", widgetHandler())
	mux.HandleFunc("/api/current", trackHandler(&t))

	go http.ListenAndServe(port, mux)
	token, auth := OAuthFlow(redirUrl)

	/* DO NOT FORGET TO CHANGE RANDSTATE LATER */

	ctx := context.Background()
	client := spotify.New(auth.Client(ctx, token))

	fmt.Println("Server is running...")

	go func() {
		for {
			playing, err := client.PlayerCurrentlyPlaying(ctx)

			if err != nil {
				fmt.Println("Error while searching for playing", err)
			}

			if playing.Item != nil {
				updateTrack(&t, playing)
			}

			//fmt.Println(&t)
			time.Sleep(4 * time.Second)
		}
	}()

	select {}

}
