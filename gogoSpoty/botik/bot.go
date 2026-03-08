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
	state := spoty.GenerateRandState()
	port := ":6111"

	clientID, ok := os.LookupEnv("CLIENT_ID")
	if !ok {
		log.Fatal("CLIENT_ID not set", "Read manual")
	}
	redirUrl, ok := os.LookupEnv("REDIRECT_URL")
	if !ok {
		log.Fatal("REDIRECT_URL not set", "Read manual")
	}

	/* FOR DEBUG ONLY */
	fmt.Println(redirUrl)
	/* ------------- */

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", callbackHandler(state))

	fmt.Println("Open url:", generateTwitchAuthUrl(clientID, redirUrl, state))

	fmt.Println("Server is running on", redirUrl)
	go http.ListenAndServe(port, mux)

	select {}

}

func generateTwitchAuthUrl(clientID string, redirUrl string, state string) string {
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

// add channeling later
func callbackHandler(checkState string, ch chan string) http.HandlerFunc {
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
