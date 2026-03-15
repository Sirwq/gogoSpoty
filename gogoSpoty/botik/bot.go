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

	port := ":6111"
	godotenv.Load(".env")

	state := spoty.GenerateRandState()
	tokenCh := make(chan string)
	username, ok := os.LookupEnv("USERNAME")
	if !ok {
		log.Fatal("USERNAME not set", "Read manual")
	}
	channelToJoin, ok := os.LookupEnv("CHANNEL")
	if !ok {
		log.Fatal("CHANNEL not set", "Read manual")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", CallbackHandler(state, tokenCh))
	go http.ListenAndServe(port, mux)

	tt, err := LoadToken()
	// Loading Token is broken if token is old

	/* TEMP CLIENT FOR TEST PURPOSES ONLY */
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

	if err != nil {

		clientID, ok := os.LookupEnv("CLIENT_ID")
		if !ok {
			log.Fatal("CLIENT_ID not set", "Read manual")
		}
		redirUrl, ok := os.LookupEnv("REDIRECT_URL")
		if !ok {
			log.Fatal("REDIRECT_URL not set", "Read manual")
		}
		clientSecret, ok := os.LookupEnv("CLIENT_SECRET")
		if !ok {
			log.Fatal("CLIENT_SECRET not set", "Read manual")
		}

		fmt.Println("Open url:", GenerateTwitchAuthUrl(clientID, redirUrl, state))
		authState := <-tokenCh

		tt, err = ExchangeCode(clientID, clientSecret, authState, redirUrl)

		if err != nil {
			panic(err)
		}

		SaveToken(tt)
	}

	clientTwitch := twitch.NewClient(username, "oauth:"+tt.AccessToken)

	requestsCooldown := NewUserCooldowns()

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
			clientTwitch.Say(channelToJoin, "You're on cooldown, wait a bit!")
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
				clientTwitch.Say(channelToJoin, "Error while adding track")
				return
			}

			answer := fmt.Sprintf("Found track: %s, Added to queue!", trackName)
			clientTwitch.Say(channelToJoin, answer)
			fmt.Println(trackID)
		}
	})

	clientTwitch.Join(channelToJoin)

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
