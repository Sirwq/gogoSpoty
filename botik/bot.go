package botik

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
}

func NewBot(
	twitchClient *twitch.Client,
	spotifyClient *spotify.Client,
	queue *Queue,
	cooldowns *UserCooldowns,
	channel string,
) *Bot {
	return &Bot{
		twitch:    twitchClient,
		spotify:   spotifyClient,
		queue:     queue,
		cooldowns: cooldowns,
		channel:   channel,
	}
}

func (bot *Bot) Start(ctx context.Context) error {
	bot.MessageHandler(ctx)
	bot.Join()
	err := bot.twitch.Connect()
	return err
}

func (bot *Bot) Join() {
	bot.twitch.Join(bot.channel)
}

func (bot *Bot) Disconnect() {
	bot.twitch.Disconnect()
}
