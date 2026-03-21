package main

import (
	"context"
	"gogoSpoty/botik"
	"gogoSpoty/spoty"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const pollInterval = 5 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	twitchConf := botik.LoadConfig()
	twitchCooldowns := botik.NewUserCooldowns()
	twitchClient, err := botik.NewTwitchClient(twitchConf, "twitchToken.json")

	if err != nil {
		log.Fatalf("Failed to create twitch client: %v", err)
	}

	redisClient := botik.NewRedisClient("1234", "localhost:6379") // change pass later
	redisQueue := botik.NewQueue(redisClient)

	spotifyConf := loadConfig()
	track := &spoty.Track{}

	spotifyClient := NewSpotifyClient(ctx, spotifyConf, "SpotifyToken.json")
	p := NewPoller(spotifyClient, track, redisQueue, 5*time.Second)

	go p.Start(ctx)

	srv := &http.Server{
		Addr:    spotifyConf.Port,
		Handler: newServer(track),
	}

	go func() {
		errCh <- srv.ListenAndServe()
	}()

	bot := botik.NewBot(
		twitchClient,
		spotifyClient,
		redisQueue,
		twitchCooldowns,
		twitchConf.TwitchChannel)

	log.Printf("Listening on: %v", spotifyConf.redirectURL)
	go bot.Start(ctx)

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), pollInterval)
	defer cancel()
	srv.Shutdown(shutdownCtx)
	redisClient.Close()
	bot.Disconnect()
	log.Println("Stopped")
}
