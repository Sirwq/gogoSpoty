package spoty_test

import (
	"encoding/json"
	"gogoSpoty/spoty"
	"net/http"
	"net/http/httptest"
	"testing"
)

var track spoty.Track

func TestTrackHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/current", nil)
	spoty.TrackHandler(&track).ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("got status code: %v, expected: 200", rec.Result().StatusCode)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("wrong content type: %v", rec.Header().Get("Content-Type"))
	}

	var result spoty.Track
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("invalid JSON: %v", err)
	}

	if result.Playing != track.Playing {
		t.Errorf("got %v, want %v", result.Playing, track.Playing)
	}
}
