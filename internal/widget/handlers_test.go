package widget

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var track Track

func TestTrackHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/current", nil)
	TrackHandler(&track).ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("got status code: %v, expected: 200", rec.Result().StatusCode)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("wrong content type: %v", rec.Header().Get("Content-Type"))
	}

	var result Track
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("invalid JSON: %v", err)
	}

	if result.Playing != track.Playing {
		t.Errorf("got %v, want %v", result.Playing, track.Playing)
	}
}

func TestWidgetHandler(t *testing.T) {
	path := "../static/widget.html"
	if _, err := os.Stat(path); err != nil {
		t.Fatal("file not found:", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/widget", nil)
	WidgetHandler("../static/widget.html").ServeHTTP(rec, req)

	location := rec.Header().Get("Location")
	t.Log("Redirect to:", location)

	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("got status code: %v, expected: 200", rec.Result().StatusCode)
	}

	var expected string = "<!DOCTYPE html>"
	body := strings.Split(rec.Body.String(), "\n")

	if body[0] != expected {
		t.Errorf("got %v, expected %v", rec.Body.String(), expected)
	}

	if body[len(body)-1] != "</html>" {
		t.Errorf("got %v, expected %v", rec.Body.String(), "</html>")
	}
}

func TestCallbackHandler(t *testing.T) {
	tokCh := make(chan *oauth2.Token)
	state := "TestingState"
	auth := spotifyauth.New()

	query := "/callback?code=ABC123&state=TestingState"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", query, nil)

	CallbackHandler(state, auth, tokCh).ServeHTTP(rec, req)

	if rec.Body == nil {
		t.Errorf("got nil, expected value")
	}

	if rec.Result().StatusCode != http.StatusForbidden {
		t.Errorf("got status code: %v, expected: 403", rec.Result().StatusCode)
	}

	body := strings.Split(rec.Body.String(), "\n")[0]

	if body != "Auth failed" {
		t.Errorf("Expected auth failure, got body:")
	}

}
