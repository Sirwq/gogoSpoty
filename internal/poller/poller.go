package poller

import (
	"context"
	"gogoSpoty/internal/bot"
	"gogoSpoty/internal/widget"
	"log"
	"time"

	"github.com/zmb3/spotify/v2"
)

type Poller struct {
	Client       *spotify.Client
	Track        *widget.Track
	Queue        *bot.Queue
	Interval     time.Duration
	LastQueuedID string
}

func NewPoller(client *spotify.Client, track *widget.Track, q *bot.Queue, interval time.Duration) *Poller {
	return &Poller{
		Client:   client,
		Track:    track,
		Queue:    q,
		Interval: interval,
	}
}

func (p *Poller) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Poller stopped")
			return
		case <-time.After(p.Interval):
			playing, err := p.Client.PlayerCurrentlyPlaying(ctx)

			if err != nil {
				log.Printf("Error fetching track %v\n", err)
				continue
			}

			if playing == nil || playing.Item == nil {
				continue
			}

			widget.UpdateTrack(p.Track, playing)
			SongRequest, err := p.Queue.Peek(ctx)

			if err == bot.ErrQueueEmpty {
				continue
			} else if err != nil {
				log.Printf("Error in poller: %v\n", err)
				continue
			}

			if playing.Item.ID == spotify.ID(SongRequest.TrackID) {
				_, err := p.Queue.Remove(ctx)
				logErr(err, "Error removing from queue in poller")
				continue
			}

			if p.LastQueuedID == SongRequest.TrackID {
				continue
			}

			err = p.Client.QueueSong(ctx, spotify.ID(SongRequest.TrackID))
			p.LastQueuedID = SongRequest.TrackID
			logErr(err, "Error queuing song in poller")
		}
	}
}

func logErr(err error, msg string) {
	if err != nil {
		log.Println(msg, "Error:"+err.Error())
	}
}
