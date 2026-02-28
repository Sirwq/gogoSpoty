package spoty_test

import (
	"gogoSpoty/spoty"
	"os"
	"testing"

	"golang.org/x/oauth2"
)

func TestAll(t *testing.T) {
	redirUrl := "123:123:123"
	r := "0.1.2.3.4"

	state, auth, ch := spoty.OAuthFlow(redirUrl)
	stateD, authD, chD := spoty.OAuthFlow(r)

	if state == "" {
		t.Errorf("state could't be empty: %v", state)
	}
	if stateD == "" {
		t.Errorf("state could't be empty: %v", stateD)
	}

	if auth == nil {
		t.Errorf("auth could't be nil: %v", auth)
	}
	if authD == nil {
		t.Errorf("auth could't be nil: %v", authD)
	}

	if ch == nil {
		t.Errorf("ch could't be nil: %v", ch)
	}
	if chD == nil {
		t.Errorf("ch could't be nil: %v", ch)
	}

	if state == stateD {
		t.Errorf("states can't be same for two authentificators: %v - %v", state, stateD)
	}

	if len(state) != 64 || len(stateD) != 64 {
		t.Errorf("State len: %v | StateD len: %v - should be 64", len(state), len(stateD))
	}

	go func() {
		ch <- &oauth2.Token{}
	}()

	emptyTok := <-ch
	var empty oauth2.Token

	if *(emptyTok) != empty {
		t.Errorf("Tokens should be equal: %v | %v", emptyTok, empty)
	}

	err := os.Setenv("CLIENT_ID", "")
	if err != nil {
		t.Fatalf("Error setting enviroment variable: %v", err)
	}
	err = os.Setenv("CLIENT_SECRET", "")
	if err != nil {
		t.Fatalf("Error setting enviroment variable: %v", err)
	}

	EmptyState, EmptyAuth, EmptyCh := spoty.OAuthFlow(r)

	if EmptyState == "" {
		t.Errorf("state could't be empty: %v", state)
	}

	if EmptyAuth == nil {
		t.Errorf("auth could't be nil: %v", state)
	}

	if EmptyCh == nil {
		t.Errorf("ch could't be nil: %v", ch)
	}

	err = os.Unsetenv("CLIENT_ID")
	if err != nil {
		t.Fatalf("Error unsetting environment variable: %v", err)
	}
	err = os.Unsetenv("CLIENT_SECRET")
	if err != nil {
		t.Fatalf("Error unsetting environment variable: %v", err)
	}

	err = spoty.SaveToken(emptyTok)

	if err != nil {
		t.Errorf("Error %v while saving token: %v", err, emptyTok)
	}

	tok, err := spoty.LoadToken()

	if err != nil {
		t.Errorf("Error %v while loading token: %v", err, tok)
	}

	t.Cleanup(func() { os.Remove("token.json") })
}
