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

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	twitchConf := botik.LoadConfig()
	twitchCooldowns := botik.NewUserCooldowns()
	twitchClient, err := botik.NewTwitchClient(twitchConf, "twitchToken.json")

	if err != nil {
		log.Fatalf("Error on twitch client creation: %v", err)
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
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(shutdownCtx)
	redisClient.Close()
	bot.Disconnect()
	log.Println("Stopped")
}
