package botik

import (
	"net/url"
	"os"
)

func main() {
	redirUrl := "https://127.0.0.1:6111/"

	r := url.Values{}
	r.Set("client_id", os.Getenv("CLIENT_ID"))
	r.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	r.Set("grant_type", "client_credentials")

	body := r.Encode()

	// https://dev.twitch.tv/docs/authentication/getting-tokens-oauth#use-the-authorization-code-to-get-a-token

}
