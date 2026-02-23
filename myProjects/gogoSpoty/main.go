package main

import (
    "context"
    "github.com/zmb3/spotify/v2"
    spotifyauth "github.com/zmb3/spotify/v2/auth"
    "golang.org/x/oauth2/clientcredentials"
)

ctx := context.Background()
config := &clientcredentials.Config{
    ClientID:     "YOUR_CLIENT_ID",
    ClientSecret: "YOUR_CLIENT_SECRET",
    TokenURL:     spotifyauth.TokenURL,
}

token, _ := config.Token(ctx)
httpClient := spotifyauth.New().Client(ctx, token)
client := spotify.New(httpClient)

func main() {
	playing, err := client.PlayerCurrentlyPlaying(ctx)
	fmt.Println(playing)
}