package app

import (
	"context"
	"gogoSpoty/internal/bot"
	"gogoSpoty/internal/config"
	"gogoSpoty/internal/poller"
	"gogoSpoty/internal/widget"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const pollInterval = 5 * time.Second

type App struct {
	srv    *http.Server
	bot    *bot.Bot
	redis  *redis.Client
	poller *poller.Poller
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

	redisClient := bot.NewRedisClient(&cfg.Redis)
	queue := bot.NewQueue(redisClient)

	b := bot.NewBot(twitch, spotify, queue, cooldowns, cfg.Twitch.Channel)
	p := poller.NewPoller(spotify, track, queue, pollInterval)
	return &App{srv: server, bot: b, redis: redisClient, poller: p}, nil
}

func (a *App) Start(ctx context.Context) {
	go a.poller.Start(ctx)
	go a.bot.Start(ctx)
	go a.srv.ListenAndServe()
}

// func (app *App) Start(ctx context.Context) {
// 	ctx, stop := signal.NotifyContext(
// 		context.Background(),
// 		syscall.SIGINT, syscall.SIGTERM)
// 	defer stop()

// 	errCh := make(chan error, 1)

//
//
// 	track := &widget.Track{}

//

// 	p :=

// 	go p.Start(ctx)

// 	go func() {
// 		errCh <- srv.ListenAndServe()
// 	}()

// 	log.Printf("Listening on: %v", &app.cfg.Spotify.RedirectURL)
// 	go bot.Start(ctx)

// 	<-ctx.Done()
// 	log.Println("Shutting down...")

// 	shutdownCtx, cancel := context.WithTimeout(context.Background(), pollInterval)
// 	defer cancel()
// 	srv.Shutdown(shutdownCtx)
// 	redis.Close()
// 	bot.Disconnect()
// 	log.Println("Stopped")
// }
