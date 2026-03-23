package app

import (
	"context"
	"fmt"
	"gogoSpoty/internal/bot"
	"gogoSpoty/internal/config"
	"gogoSpoty/internal/poller"
	"gogoSpoty/internal/widget"
	"io"
	"log"
	"net/http"
	"time"
)

const pollInterval = 5 * time.Second

type App struct {
	srv    *http.Server
	bot    *bot.Bot
	closer io.Closer
	poller *poller.Poller
	errors chan error
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {

	track := &widget.Track{}

	server := &http.Server{
		Addr:    cfg.Spotify.Port,
		Handler: widget.NewServerMux(track),
	}

	cooldowns := bot.NewUserCooldowns()
	twitch, err := bot.NewTwitchClient(&cfg.Twitch, "TwitchToken.json")

	if err != nil {
		return nil, err
	}

	spotify := widget.NewSpotifyClient(ctx, &cfg.Spotify, "SpotifyToken.json")
	queue := bot.NewQueue(&cfg.Redis)

	b := bot.NewBot(twitch, spotify, queue, cooldowns, cfg.Twitch.Channel)
	p := poller.NewPoller(spotify, track, queue, pollInterval)
	e := make(chan error, 2)
	return &App{srv: server, bot: b, closer: queue, poller: p, errors: e}, nil
}

func (a *App) Start(ctx context.Context) {
	go a.poller.Start(ctx)
	go func() {
		if err := a.bot.Start(ctx); err != nil {
			a.errors <- fmt.Errorf("bot: %w", err)
		}
	}()
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.errors <- fmt.Errorf("http server: %w", err)
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) {
	if err := a.srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
	if a.closer != nil {
		if err := a.closer.Close(); err != nil {
			log.Printf("Redis shutdown error: %v", err)
		}
	}
	a.bot.Disconnect()
}

func (a *App) Errors() <-chan error {
	return a.errors
}
