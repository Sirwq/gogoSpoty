package main

import (
	"context"
	"fmt"
	"gogoSpoty/helpers"
	"gogoSpoty/spoty"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

type spotifyConfig struct {
	clientID     string
	clientSecret string
	redirectURL  string
	Port         string
}

func NewSpotifyClient(ctx context.Context, config *spotifyConfig, tokName string) *spotify.Client {
	state, auth, ch := spoty.OAuthFlow(config.redirectURL, config.clientID, config.clientSecret)
	token, err := spoty.LoadToken(tokName)

	if err != nil {
		mux := http.NewServeMux()
		mux.HandleFunc("/callback", spoty.CallbackHandler(state, auth, ch))
		srv := &http.Server{Addr: config.Port, Handler: mux}

		go srv.ListenAndServe()

		url := auth.AuthURL(state)
		fmt.Println("Open this url: ", url)
		token = <-ch
		srv.Shutdown(ctx)
		spoty.SaveToken(token, tokName)
	}

	return spotify.New(auth.Client(ctx, token))
}

func loadConfig() *spotifyConfig {
	var ok bool
	var conf spotifyConfig
	godotenv.Load(".env")

	conf.clientID, ok = os.LookupEnv("CLIENT_ID_SPOTY")
	helpers.CheckErrFatal(ok, "CLIENT_ID_SPOTY not set")

	conf.clientSecret, ok = os.LookupEnv("CLIENT_SECRET_SPOTY")
	helpers.CheckErrFatal(ok, "CLIENT_SECRET_SPOTY not set")

	conf.redirectURL, ok = os.LookupEnv("REDIRECT_URL_SPOTY")
	helpers.CheckErrFatal(ok, "REDIRECT_URL_SPOTY not set")

	conf.Port = ":5111"
	return &conf
}
