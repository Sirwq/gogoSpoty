package main

import (
	"context"
	"gogoSpoty/botik"
	"gogoSpoty/helpers"
	"gogoSpoty/spoty"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	twitchConf := botik.LoadConfig()
	twitchCooldowns := botik.NewUserCooldowns()
	twitchClient, err := botik.NewTwitchClient(twitchConf, "twitchToken.json")

	if err != nil {
		helpers.CheckErrFatal(false, "Error on twitch client creation")
		return
	}

	redisClient := botik.NewRedisClient("1234", "localhost:6379") // change pass later
	redisQueue := botik.NewQueue(redisClient)

	var duration time.Duration = 5

	spotifyConf := loadConfig()
	track := &spoty.Track{}
	mux := newServer(track)

	spotifyClient := NewSpotifyClient(ctx, spotifyConf, "SpotifyToken.json")
	p := NewPoller(spotifyClient, track, redisQueue, duration)

	go p.Start(ctx)
	go http.ListenAndServe(spotifyConf.Port, mux)

	bot := botik.NewBot(
		twitchClient,
		spotifyClient,
		redisQueue,
		twitchCooldowns,
		twitchConf.TwitchChannel)

	bot.Start(ctx)
	helpers.Log("Bot is working...")
}
