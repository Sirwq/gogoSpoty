package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type twitchToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (tt twitchToken) String() string {
	return fmt.Sprintf("Access Token: %s\nRefresh Token: %s, Expiry ms: %d", tt.AccessToken, tt.RefreshToken, tt.ExpiresIn)
}

func GenerateTwitchAuthUrl(clientID string, redirUrl string, state string) string {
	data := url.Values{}

	data.Set("client_id", clientID)
	data.Set("redirect_uri", redirUrl)
	data.Set("response_type", "code")
	data.Set("scope", "chat:read chat:edit")
	data.Set("state", state)

	encodedQuery := data.Encode()
	authUrl := "https://id.twitch.tv/oauth2/authorize?" + encodedQuery

	return authUrl
}

func ExchangeCode(clientID string, clientSecret string, code string, redirUrl string) (*twitchToken, error) {
	data := url.Values{}

	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirUrl)

	encodedData := data.Encode()
	url := "https://id.twitch.tv/oauth2/token"

	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(encodedData))
	var tok twitchToken

	if err != nil {
		fmt.Printf("Got error %v\n", err)
		return &tok, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tok)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("error decoding json: %v", err)
			return &tok, err
		}
	}

	return &tok, nil
}

func CallbackHandler(checkState string, ch chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		err := r.URL.Query().Get("error")

		if err != "" {
			http.Error(w, "Auth failed", http.StatusBadRequest)
			return
		}

		if code == "" {
			http.Error(w, "code not found", http.StatusBadRequest)
			return
		}

		if state != checkState {
			http.Error(w, "state is falcificated", http.StatusExpectationFailed)
			return
		}

		ch <- code
		fmt.Fprintf(w, "Autorization code recieved!")
	}
}

func SaveToken(tt *twitchToken) error {
	data, err := json.Marshal(tt)

	if err != nil {
		return err
	}
	return os.WriteFile("token.json", data, 0600)
}

func LoadToken() (*twitchToken, error) {
	data, err := os.ReadFile("token.json")
	if err != nil {
		return nil, err
	}
	var token twitchToken
	err = json.Unmarshal(data, &token)
	return &token, err
}
