package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func main() {
	godotenv.Load(".env")

	port := ":5111"

	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv("REDIRECT_URI")),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
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

	go http.ListenAndServe(port, nil) // creating a server without handler

	url := auth.AuthURL(state)
	fmt.Println("Open this url:", url)

	token := <-ch
	client := spotify.New(auth.Client(context.Background(), token))

	user, _ := client.CurrentUser(context.Background())
	fmt.Println("Logged in as:", user.DisplayName)

}
