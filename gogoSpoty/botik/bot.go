package main

import (
	"context"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/zmb3/spotify/v2"
)

type Bot struct {
	twitch    *twitch.Client
	spotify   *spotify.Client
	queue     *Queue
	cooldowns *UserCooldowns
	channel   string
	ctx       context.Context
}

func NewBot(
	twitchClient *twitch.Client,
	spotifyClient *spotify.Client,
	queue *Queue,
	cooldowns *UserCooldowns,
	channel string,
	ctx context.Context,
) *Bot {
	return &Bot{
		twitch:    twitchClient,
		spotify:   spotifyClient,
		queue:     queue,
		cooldowns: cooldowns,
		channel:   channel,
		ctx:       ctx,
	}
}

func (bot *Bot) Start() error {
	bot.Join()
	err := bot.twitch.Connect()
	bot.MessageHandler()
	return err
}

func (bot *Bot) Join() {
	bot.twitch.Join(bot.channel)
}

// 	twitchConf := LoadConfig()
// 	requestsCooldown := NewUserCooldowns()
// 	clientTwitch, err := NewTwitchClient(twitchConf, "twitchToken.json")

// 	if err != nil {
// 		fmt.Println("error on twitch client creation: ", err)
// 		return
// 	}

// }
