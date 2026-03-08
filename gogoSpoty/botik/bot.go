package main

import (
	"fmt"
	"gogoSpoty/spoty"
	"log"
	"net/http"
	"os"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
)

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

	client := twitch.NewClient(username, "oauth:"+tt.AccessToken)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)
	})

	client.Join(channelToJoin)

	fmt.Println("Bot is running")
	err = client.Connect()
	if err != nil {
		panic(err)
	}

}
