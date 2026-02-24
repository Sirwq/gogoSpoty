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

	go http.ListenAndServe(port, nil)
	token, auth := OAuthFlow(redirUrl)

	/* DO NOT FORGET TO CHANGE RANDSTATE LATER */

	ctx := context.Background()
	client := spotify.New(auth.Client(ctx, token))

	go func() {
		var t Track
		for {
			playing, err := client.PlayerCurrentlyPlaying(ctx)

			if err != nil {
				fmt.Println("Error while searching for playing", err)
			}

			if playing.Item != nil {
				updateTrack(&t, playing)
			}

			fmt.Println(&t)
			time.Sleep(5 * time.Second)
		}
	}()

	select {}

}
