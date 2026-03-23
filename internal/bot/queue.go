package bot

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrQueueEmpty = errors.New("queue is empty")

type Queue interface {
	Add(ctx context.Context, req SongRequest) error
	Remove(ctx context.Context) (SongRequest, error)
	Peek(ctx context.Context) (SongRequest, error)
	Read(ctx context.Context) ([]SongRequest, error)
	Close() error
}

type SongRequest struct {
	Username    string    `json:"username"`
	TrackID     string    `json:"track_id"`
	TrackName   string    `json:"track_name"`
	TrackArtist string    `json:"track_artist"`
	RequestedAt time.Time `json:"requested_at"`
}

func (s *SongRequest) DisplayName() string {
	return fmt.Sprintf("%s - %s", s.TrackName, s.TrackArtist)
}
