package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/zmb3/spotify/v2"
)

const PREFIX_REQUEST = "!sr"

func (bot *Bot) MessageHandler() {
	bot.twitch.OnPrivateMessage(func(message twitch.PrivateMessage) {
		m := message.Message
		uname := message.User.Name

		m, ok := SongQuery(m, PREFIX_REQUEST)
		if !ok {
			return
		}

		if m == "" {
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
		results, err := bot.spotify.Search(bot.ctx, m, spotify.SearchTypeTrack)

		if err != nil {
			fmt.Println(err)
			return
		}

		if len(results.Tracks.Tracks) == 0 {
			fmt.Println("Track not found") // Answer in caht later
		} else {
			fmt.Println(time.Since(lastRequest))

			reqTime := time.Now()
			bot.cooldowns.Store(uname, reqTime)
			// trackID := r.Tracks.Tracks[0].ID
			// trackName := r.Tracks.Tracks[0].Name
			// err := bot.spotify.QueueSong(bot.ctx, trackID)
			var req SongRequest = SongRequest{
				Usename:     uname,
				RequestedAt: reqTime,
				TrackID:     results.Tracks.Tracks[0].ID,
				TrackName:   results.Tracks.Tracks[0].Name,
				TrackArtist: results.Tracks.Tracks[0].Artists,
			}

			if err != nil {
				bot.twitch.Say(bot.channel, "Error while adding track")
				return
			}

			answer := fmt.Sprintf("Found track: %s, Added to queue!", req.TrackName)
			bot.twitch.Say(bot.channel, answer)
			fmt.Println(req.TrackID)
		}
	})

}

func SongQuery(msg, prefix string) (string, bool) {
	_, after, found := strings.Cut(msg, prefix)

	if !found {
		return "", false
	}

	return strings.TrimSpace(after), true

}
