package widget

import (
	"context"
	"fmt"
	"gogoSpoty/internal/config"
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

		go srv.ListenAndServe()

		url := auth.AuthURL(state)
		fmt.Println("Open this url: ", url)
		token = <-ch
		srv.Shutdown(ctx)
		SaveToken(token, tokName)
	}

	return spotify.New(auth.Client(ctx, token))
}
