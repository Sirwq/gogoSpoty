package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Twitch  TwitchConfig
	Spotify SpotifyConfig
	Redis   RedisConfig
}

type TwitchConfig struct {
	Username     string
	Channel      string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Port         string
}

type SpotifyConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Port         string
}

type RedisConfig struct {
	Addr     string
	Password string
}

func Load() *Config {
	godotenv.Load(".env")

	return &Config{
		Twitch:  loadTwitch(),
		Spotify: loadSpotify(),
		Redis:   loadRedis(),
	}

}

func loadTwitch() TwitchConfig {
	var ok bool
	var cfg TwitchConfig
	cfg.Port = ":6111"
	cfg.Username, ok = os.LookupEnv("TWITCH_USERNAME")
	check(ok, "twitch username not set")

	cfg.Channel, ok = os.LookupEnv("TWITCH_CHANNEL")
	check(ok, "twitch channel not set")

	cfg.ClientID, ok = os.LookupEnv("TWITCH_CLIENT_ID")
	check(ok, "twitch clientID not set")

	cfg.RedirectURL, ok = os.LookupEnv("TWITCH_REDIRECT_URL")
	check(ok, "twitch redirect url not set")

	cfg.ClientSecret, ok = os.LookupEnv("TWITCH_CLIENT_SECRET")
	check(ok, "twitch client secret not set")
	return cfg
}

func loadSpotify() SpotifyConfig {
	var ok bool
	var cfg SpotifyConfig
	cfg.Port = ":5111"

	cfg.ClientID, ok = os.LookupEnv("CLIENT_ID_SPOTY")
	check(ok, "CLIENT_ID_SPOTY not set")

	cfg.ClientSecret, ok = os.LookupEnv("CLIENT_SECRET_SPOTY")
	check(ok, "CLIENT_SECRET_SPOTY not set")

	cfg.RedirectURL, ok = os.LookupEnv("REDIRECT_URL_SPOTY")
	check(ok, "REDIRECT_URL_SPOTY not set")
	return cfg
}

func loadRedis() RedisConfig {
	var ok bool
	var cfg RedisConfig
	cfg.Addr, ok = os.LookupEnv("REDIS_ADDR")
	check(ok, "REDIS_ADDR not set")

	cfg.Password, ok = os.LookupEnv("REDIS_PASSWORD")
	check(ok, "REDIS_PASSWORD not set")
	return cfg
}

func check(ok bool, msg string) {
	if !ok {
		log.Fatal(msg)
	}
}
