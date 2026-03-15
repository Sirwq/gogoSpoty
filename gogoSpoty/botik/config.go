package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type TwitchConfig struct {
	TwitchUsername     string
	TwitchChannel      string
	TwitchClientID     string
	TwitchClientSecret string
	TwitchRedirectURL  string
}

func LoadConfig() *TwitchConfig {
	var ok bool
	var conf TwitchConfig
	godotenv.Load(".env")

	conf.TwitchUsername, ok = os.LookupEnv("USERNAME")
	checkErr(ok, "twitch username not set")

	conf.TwitchChannel, ok = os.LookupEnv("CHANNEL")
	checkErr(ok, "twitch channel not set")

	conf.TwitchClientID, ok = os.LookupEnv("CLIENT_ID")
	checkErr(ok, "twitch clientID not set")

	conf.TwitchRedirectURL, ok = os.LookupEnv("REDIRECT_URL")
	checkErr(ok, "twitch redirect url not set")

	conf.TwitchClientSecret, ok = os.LookupEnv("CLIENT_SECRET")
	checkErr(ok, "twitch client secret not set")

	return &conf
}

func checkErr(ok bool, msg string) {
	if !ok {
		log.Fatal(msg, "\nRead manual")
	}
}
