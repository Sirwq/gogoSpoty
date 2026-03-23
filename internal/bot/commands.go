package bot

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/zmb3/spotify/v2"
)

const prefixRequest = "!sr"

func (bot *Bot) MessageHandler(ctx context.Context) {
	bot.twitch.OnPrivateMessage(func(message twitch.PrivateMessage) {
		m := message.Message
		uname := message.User.Name

		m, ok := SongQuery(m, prefixRequest)
		if !ok || m == "" {
			return
		}

		lastRequest, ok := bot.cooldowns.Load(uname)

		if !ok {
			bot.cooldowns.Store(uname, time.Time{})
			lastRequest = time.Time{}
		}

		if time.Since(lastRequest) < cooldown*time.Second {
			bot.twitch.Say(bot.channel, "You're on cooldown, wait a bit!")
			return
		}

		// Search request. Assumes that first entry is that what user wished
		results, err := bot.spotify.Search(ctx, m, spotify.SearchTypeTrack)

		if err != nil {
			log.Println(err)
			return
		}

		if len(results.Tracks.Tracks) == 0 {
			bot.twitch.Say(bot.channel, "Track not found")
			return
		}

		reqTime := time.Now()
		req := SongRequest{
			Username:    uname,
			RequestedAt: reqTime,
			TrackID:     string(results.Tracks.Tracks[0].ID),
			TrackName:   results.Tracks.Tracks[0].Name,
		}

		var artists []string
		if results.Tracks.Tracks[0].Artists != nil {
			for _, artist := range results.Tracks.Tracks[0].Artists {
				artists = append(artists, artist.Name)
			}
		}

		req.TrackArtist = strings.Join(artists, ", ")
		err = bot.queue.Add(ctx, req)

		var answer string

		if err != nil {
			log.Println("Failed to push to queue")
			answer = "failed to push to queue"
		} else {
			bot.cooldowns.Store(uname, reqTime)
			answer = ("found: " + req.DisplayName() + ", added to queue")
		}

		bot.twitch.Say(bot.channel, answer)

	})

}

func SongQuery(msg, prefix string) (string, bool) {
	if !strings.HasPrefix(msg, prefix) {
		return "", false
	}

	return strings.TrimSpace(msg[len(prefix):]), true
}
