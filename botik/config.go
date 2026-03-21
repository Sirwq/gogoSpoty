package botik

import (
	"gogoSpoty/helpers"
	"os"

	"github.com/joho/godotenv"
)

type TwitchConfig struct {
	TwitchUsername     string
	TwitchChannel      string
	TwitchClientID     string
	TwitchClientSecret string
	TwitchRedirectURL  string
	TwitchPort         string
}

func LoadConfig() *TwitchConfig {
	var ok bool
	var conf TwitchConfig
	godotenv.Load(".env")

	conf.TwitchUsername, ok = os.LookupEnv("TWITCH_USERNAME")
	helpers.CheckErrFatal(ok, "twitch username not set")

	conf.TwitchChannel, ok = os.LookupEnv("TWITCH_CHANNEL")
	helpers.CheckErrFatal(ok, "twitch channel not set")

	conf.TwitchClientID, ok = os.LookupEnv("TWITCH_CLIENT_ID")
	helpers.CheckErrFatal(ok, "twitch clientID not set")

	conf.TwitchRedirectURL, ok = os.LookupEnv("TWITCH_REDIRECT_URL")
	helpers.CheckErrFatal(ok, "twitch redirect url not set")

	conf.TwitchClientSecret, ok = os.LookupEnv("TWITCH_CLIENT_SECRET")
	helpers.CheckErrFatal(ok, "twitch client secret not set")

	conf.TwitchPort = ":6111"

	return &conf
}
