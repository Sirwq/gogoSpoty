package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type Track struct {
	mx        sync.Mutex
	Item      spotify.FullTrack       `json:"item"`
	Playing   bool                    `json:"is_playing"`
	Progress  spotify.Numeric         `json:"progress_ms"`
	Timestamp int64                   `json:"timestamp"`
	Context   spotify.PlaybackContext `json:"context"`
}

func (t *Track) String() string {
	var s strings.Builder
	t.mx.Lock()
	s.WriteString(t.Item.Name + "\n")
	for i, v := range t.Item.Artists {
		s.WriteString(v.Name)
		if i < len(t.Item.Artists)-1 {
			s.WriteString(", ")
		}
	}
	t.mx.Unlock()
	return s.String()
}

func updateTrack(t *Track, playing *spotify.CurrentlyPlaying) {
	t.mx.Lock()
	if playing.Playing {
		t.Item = *playing.Item
		t.Playing = playing.Playing
		t.Timestamp = playing.Timestamp
		t.Context = playing.PlaybackContext
	} else {
		fmt.Println("Paused or nothing is playing")
	}
	t.mx.Unlock()
}

func main() {
	godotenv.Load(".env")

	port := ":5111"
	redirUrl := "http://127.0.0.1:5111/callback"

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

	go http.ListenAndServe(port, nil) // creating a server without handler

	token, err := loadToken()
	if err != nil {
		url := auth.AuthURL(state)
		fmt.Println("Open this url: ", url)
		token = <-ch
		saveToken(token)
	}

	client := spotify.New(auth.Client(context.Background(), token))
	ctx := context.Background()

	go func() {
		var t Track
		for {
			playing, _ := client.PlayerCurrentlyPlaying(ctx)

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
