package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func OAuthFlow(redirUrl string) (*oauth2.Token, *spotifyauth.Authenticator) {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirUrl),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
		),
		spotifyauth.WithClientID(os.Getenv("CLIENT_ID")),
		spotifyauth.WithClientSecret(os.Getenv("CLIENT_SECRET")),
	)

	state := "randstate"

	ch := make(chan *oauth2.Token)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Auth failed", http.StatusForbidden)
			return
		}
		fmt.Fprintf(w, "Login success!")
		ch <- token
	})

	token, err := loadToken()
	if err != nil {
		url := auth.AuthURL(state)
		fmt.Println("Open this url: ", url)
		token = <-ch
		saveToken(token)
	}
	return token, auth
}

func saveToken(token *oauth2.Token) error {
	data, err := json.Marshal(token)

	if err != nil {
		return err
	}
	return os.WriteFile("token.json", data, 0600)
}

func loadToken() (*oauth2.Token, error) {
	data, err := os.ReadFile("token.json")
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal(data, &token)
	return &token, err
}
