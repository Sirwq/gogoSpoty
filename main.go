package main

import (
	"context"
	"fmt"
	"gogoSpoty/botik"
	"gogoSpoty/spoty"
	"net/http"
	"time"
)

/* implement .ENV loading trough UI */

func main() {
	ctx := context.Background()

	twitchConf := botik.LoadConfig()
	twitchCooldowns := botik.NewUserCooldowns()
	twitchClient, err := botik.NewTwitchClient(twitchConf, "twitchToken.json")

	if err != nil {
		fmt.Println("error on twitch client creation: ", err)
		return
	}

	redisClient := botik.NewRedisClient("1234") // change pass later
	redisQueue := botik.NewQueue(redisClient)

	port := ":5111"
	var duration time.Duration = 5

	spotifyConf := loadConfig()
	track := &spoty.Track{}
	mux := newServer(track)
	spotifyClient := NewSpotifyClient(ctx, spotifyConf, mux, "SpotifyToken.json")

	p := &Poller{
		Client:   spotifyClient,
		Track:    track,
		Queue:    redisQueue,
		Interval: duration,
	}

	go p.Start(ctx)
	go http.ListenAndServe(port, mux)
	bot := botik.NewBot(
		twitchClient,
		spotifyClient,
		redisQueue,
		twitchCooldowns,
		twitchConf.TwitchChannel)

	bot.Start(ctx)
	bot.Join()
}
