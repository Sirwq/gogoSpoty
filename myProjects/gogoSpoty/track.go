package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/zmb3/spotify/v2"
)

type Track struct {
	mx        sync.Mutex
	Item      spotify.FullTrack       `json:"item"`
	Playing   bool                    `json:"is_playing"`
	Progress  spotify.Numeric         `json:"progress_ms"`
	Timestamp int64                   `json:"timestamp"`
	Context   spotify.PlaybackContext `json:"context"`
}

func (t *Track) String() string {
	var s strings.Builder
	t.mx.Lock()
	s.WriteString(t.Item.Name + "\n")
	for i, v := range t.Item.Artists {
		s.WriteString(v.Name)
		if i < len(t.Item.Artists)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString("\n")
	s.WriteString(t.Item.Album.Images[0].URL)
	t.mx.Unlock()
	return s.String()
}

func updateTrack(t *Track, playing *spotify.CurrentlyPlaying) {
	t.mx.Lock()
	if playing.Playing {
		t.Item = *playing.Item
		t.Playing = playing.Playing
		t.Timestamp = playing.Timestamp
		t.Context = playing.PlaybackContext
	} else {
		fmt.Println("Paused or nothing is playing")
	}
	t.mx.Unlock()
}
