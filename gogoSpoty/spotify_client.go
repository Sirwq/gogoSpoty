package main

import (
	"context"
	"fmt"
	"gogoSpoty/spoty"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

type spotifyConfig struct {
	clientID     string
	clientSecret string
}

func NewSpotifyClient(ctx context.Context, config *spotifyConfig, mux *http.ServeMux) *spotify.Client {

	redirUrl := "http://127.0.0.1:5111/callback"

	state, auth, ch := spoty.OAuthFlow(redirUrl, config.clientID, config.clientSecret)
	token, err := spoty.LoadToken("token.json")

	if err != nil {
		http.HandleFunc("/callback", spoty.CallbackHandler(state, auth, ch))
		go http.ListenAndServe(":5111", mux)
		url := auth.AuthURL(state)
		fmt.Println("Open this url: ", url)
		token = <-ch
		spoty.SaveToken(token, "token.json")
	}

	return spotify.New(auth.Client(ctx, token))
}

func loadConfig() *spotifyConfig {
	var ok bool
	var conf spotifyConfig
	godotenv.Load(".env")

	conf.clientID, ok = os.LookupEnv("CLIENT_ID_SPOTY")
	checkErr(ok, "CLIENT_ID_SPOTY not set")

	conf.clientSecret, ok = os.LookupEnv("CLIENT_SECRET_SPOTY")
	checkErr(ok, "Spotify CLIENT_SECRET not set")
	return &conf
}

func checkErr(ok bool, msg string) {
	if !ok {
		log.Fatal(msg, "\nRead manual")
	}
}
