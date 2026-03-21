package app

import (
	"context"
	"gogoSpoty/internal/bot"
	"gogoSpoty/internal/config"
	"gogoSpoty/internal/widget"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const pollInterval = 5 * time.Second

type App struct {
	srv   *http.Server
	bot   *bot.Bot
	redis *redis.Client
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {

	mux := widget.NewServerMux(&widget.Track{})

	server := http.ListenAndServe("", mux)

	cooldowns := bot.NewUserCooldowns()
	twitch, err := bot.NewTwitchClient(&cfg.Twitch, "TwitchToken.json")
	if err != nil {
		log.Fatalf("Failed to create twitch client: %v", err)
	}

	spotify := widget.NewSpotifyClient(ctx, &cfg.Spotify, "SpotifyToken.json")

	redisClient := bot.NewRedisClient(&cfg.Redis)
	queue := bot.NewQueue(redisClient)

	b := bot.NewBot(twitch, spotify, queue, cooldowns, cfg.Twitch.Channel)
	return &App{srv: server, bot: b, redis: redisClient}, nil
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

// 	p := poller.NewPoller(spotify, track, queue, pollInterval)

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
