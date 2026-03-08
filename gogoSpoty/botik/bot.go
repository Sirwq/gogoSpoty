package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

func main() {
	//redirUrl := "https://127.0.0.1:6111/"
	//state := "RANDSTATE" // change later
	username := "vvxshin"
	oauth, err := GetTwitchToken()

	if err != nil {
		panic(err)
	}

	client := twitch.NewClient(username, oauth)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)
	})

	client.Join("gogoBotik")

	err = client.Connect()

	if err != nil {
		panic(err)
	}
	fmt.Println(oauth)
	fmt.Println("Success!")

}

func GetTwitchToken() (string, error) {
	r := url.Values{}
	r.Set("client_id", os.Getenv("CLIENT_ID"))
	r.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	r.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(r.Encode()))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form- urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.AccessToken, nil
}
