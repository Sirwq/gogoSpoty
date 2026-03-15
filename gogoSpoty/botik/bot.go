package main

import (
	"context"
	"fmt"
	"gogoSpoty/spoty"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

type UserCooldowns struct {
	sync.Mutex
	m map[string]time.Time
}

func NewUserCooldowns() *UserCooldowns {
	return &UserCooldowns{
		m: map[string]time.Time{},
	}
}

func (uc *UserCooldowns) Store(key string, time time.Time) {
	uc.Mutex.Lock()
	uc.m[key] = time
	uc.Mutex.Unlock()
}

func (uc *UserCooldowns) Load(key string) (time.Time, bool) {
	uc.Mutex.Lock()
	v, ok := uc.m[key]
	uc.Mutex.Unlock()
	return v, ok
}

const prefSongReques = "!sr"
const cooldown = 10

func main() {
	/* TEMP CLIENT FOR TEST PURPOSES ONLY */
	godotenv.Load(".env")
	clientIDspoty, ok := os.LookupEnv("CLIENT_ID_SPOTY")
	if !ok {
		log.Fatal("Spotify CLIENT_ID_SPOTY not set, read manual")
	}

	clientSecretSpoty, ok := os.LookupEnv("CLIENT_SECRET_SPOTY")
	if !ok {
		log.Fatal("Spotify CLIENT_SECRET_SPOTY not set", "Read manual")
	}

	redirUrl := "http://127.0.0.1:5111/callback"

	stateSpoty, auth, ch := spoty.OAuthFlow(redirUrl, clientIDspoty, clientSecretSpoty)
	muxS := http.NewServeMux()
	muxS.HandleFunc("/callback", spoty.CallbackHandler(stateSpoty, auth, ch))

	go http.ListenAndServe(":5111", muxS)
	ctx := context.Background()
	token, err := spoty.LoadToken("Spoty.json")
	if err != nil {
		url := auth.AuthURL(stateSpoty)
		fmt.Println("Open this url for spoty auth: ", url)
		token = <-ch
		spoty.SaveToken(token, "Spoty.json")
	}

	clientSpotify := spotify.New(auth.Client(ctx, token))

	/* TEMP CLIENT FOR TEST PURPOSES ONLY */

	requestsCooldown := NewUserCooldowns()

	twitchConf := LoadConfig()

	clientTwitch, err := NewTwitchClient(twitchConf, "twitchToken.json")

	if err != nil {
		fmt.Println("error on twitch client creation: ", err)
		return
	}

	clientTwitch.OnPrivateMessage(func(message twitch.PrivateMessage) {
		m := message.Message
		uname := message.User.Name

		m, ok = parseCommand(m, prefSongReques)
		if !ok {
			return
		}

		if m == "" {
			return
		}

		lastRequest, ok := requestsCooldown.Load(uname)

		if !ok {
			requestsCooldown.Store(uname, time.Time{})
			lastRequest = time.Time{}
		}

		if time.Since(lastRequest) < cooldown*time.Second {
			clientTwitch.Say(twitchConf.TwitchChannel, "You're on cooldown, wait a bit!")
			return
		}

		r, err := clientSpotify.Search(ctx, m, spotify.SearchTypeTrack)

		if err != nil {
			fmt.Println(err)
			return
		}

		if len(r.Tracks.Tracks) == 0 {
			fmt.Println("Track not found") // Answer in caht later
		} else {
			fmt.Println(time.Since(lastRequest))

			requestsCooldown.Store(uname, time.Now())
			trackID := r.Tracks.Tracks[0].ID
			trackName := r.Tracks.Tracks[0].Name
			err := clientSpotify.QueueSong(ctx, trackID)

			if err != nil {
				clientTwitch.Say(twitchConf.TwitchChannel, "Error while adding track")
				return
			}

			answer := fmt.Sprintf("Found track: %s, Added to queue!", trackName)
			clientTwitch.Say(twitchConf.TwitchChannel, answer)
			fmt.Println(trackID)
		}
	})

	clientTwitch.Join(twitchConf.TwitchChannel)

	fmt.Println("Bot is running")
	err = clientTwitch.Connect()
	if err != nil {
		panic(err)
	}

}

func parseCommand(msg, prefix string) (string, bool) {
	_, after, found := strings.Cut(msg, prefix)

	if !found {
		return "", false
	}

	return strings.TrimSpace(after), true

}
