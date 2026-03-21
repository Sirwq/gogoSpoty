package widget

import (
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

/*  TODO add accent color to Track struct that will help with colors in js*/

func (t *Track) String() string {
	t.mx.Lock()
	defer t.mx.Unlock()

	var s strings.Builder
	s.WriteString(t.Item.Name + "\n")
	for i, v := range t.Item.Artists {
		s.WriteString(v.Name)
		if i < len(t.Item.Artists)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString("\n")
	return s.String()
}

func UpdateTrack(t *Track, playing *spotify.CurrentlyPlaying) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.Item = *playing.Item
	t.Playing = playing.Playing
	t.Timestamp = playing.Timestamp
	t.Progress = playing.Progress
	t.Context = playing.PlaybackContext
}
