package main

import (
	"fmt"
	"gogoSpoty/spoty"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	//port := ":6111"

	clientID, ok := os.LookupEnv("CLIENT_ID")
	if !ok {
		log.Fatal("CLIEND_ID not set", "Read manual")
	}
	redirUrl, ok := os.LookupEnv("REDIRECT_URL")
	if !ok {
		log.Fatal("REDIRECT_URL not set", "Read manual")
	}

	redirUrl += "callback"

	fmt.Println(redirUrl)

	mux := http.NewServeMux()
	go http.ListenAndServe(redirUrl+"/callback", mux)

	generateTwitchAuthUrl(clientID, redirUrl)
}

func generateTwitchAuthUrl(clientID string, redirUrl string) {
	state := spoty.GenerateRandState()
	data := url.Values{}

	data.Set("client_id", clientID)
	data.Set("redirect_uri", redirUrl)
	data.Set("response_type", "code")
	data.Set("scope", "chat:read chat:edit")
	data.Set("state", state)

	encodedQuery := data.Encode()
	authUrl := "https://id.twitch.tv/oauth2/authorize?" + encodedQuery

	fmt.Println("Open this url: ", authUrl)
}
