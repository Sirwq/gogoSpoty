package widget

import (
	"context"
	"fmt"
	"gogoSpoty/internal/config"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
)

func NewSpotifyClient(ctx context.Context, config *config.SpotifyConfig, tokName string) *spotify.Client {
	state, auth, ch := OAuthFlow(config.RedirectURL, config.ClientID, config.ClientSecret)
	token, err := LoadToken(tokName)

	if err != nil {
		mux := http.NewServeMux()
		mux.HandleFunc("/callback", CallbackHandler(state, auth, ch))
		srv := &http.Server{Addr: config.Port, Handler: mux}

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("Failed to create temp server: %v\n", err)
			}
		}()

		url := auth.AuthURL(state)
		fmt.Println("Open this url: ", url)
		token = <-ch
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown temp server: %v\n", err)
		}
		err = SaveToken(token, tokName)
		if err != nil {
			log.Printf("Failed to save Spotify token: %v\n", err)
		}
	}

	return spotify.New(auth.Client(ctx, token))
}
