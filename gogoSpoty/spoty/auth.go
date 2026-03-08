package spoty

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func OAuthFlow(redirUrl string) (string, *spotifyauth.Authenticator, chan *oauth2.Token) {
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

	state := GenerateRandState()
	ch := make(chan *oauth2.Token)
	return state, auth, ch
}

func SaveToken(token *oauth2.Token) error {
	data, err := json.Marshal(token)

	if err != nil {
		return err
	}
	return os.WriteFile("token.json", data, 0600)
}

func LoadToken() (*oauth2.Token, error) {
	data, err := os.ReadFile("token.json")
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal(data, &token)
	return &token, err
}

func GenerateRandState() string {
	k := make([]byte, 32)
	rand.Read(k)
	return hex.EncodeToString(k)
}
