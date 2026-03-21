package main

import (
	"context"
	"fmt"
	"gogoSpoty/botik"
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
		fmt.Println("error on twitch client creation: ", err)
		return
	}

	redisClient := botik.NewRedisClient("1234", "localhost:6379") // change pass later
	redisQueue := botik.NewQueue(redisClient)

	var duration time.Duration = 5

	spotifyConf := loadConfig()
	track := &spoty.Track{}
	mux := newServer(track)

	spotifyClient := NewSpotifyClient(ctx, spotifyConf, "SpotifyToken.json")

	p := &Poller{
		Client:   spotifyClient,
		Track:    track,
		Queue:    redisQueue,
		Interval: duration,
	}

	go p.Start(ctx)
	go http.ListenAndServe(spotifyConf.Port, mux)

	bot := botik.NewBot(
		twitchClient,
		spotifyClient,
		redisQueue,
		twitchCooldowns,
		twitchConf.TwitchChannel)

	bot.Start(ctx)
}
