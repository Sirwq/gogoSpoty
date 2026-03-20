package main

import (
	"context"
	"gogoSpoty/spoty"
	"log"
	"time"

	"github.com/zmb3/spotify/v2"
)

type Poller struct {
	Client   *spotify.Client
	Track    *spoty.Track
	Interval time.Duration
}

func (p *Poller) Start(ctx context.Context) {
	go func() {
		for {
			playing, err := p.Client.PlayerCurrentlyPlaying(ctx)

			if err != nil {
				log.Printf("Error fetching track %v", err)
				time.Sleep(p.Interval * time.Second)
				continue
			}

			if playing != nil && playing.Item != nil {
				spoty.UpdateTrack(p.Track, playing)
			}

			time.Sleep(p.Interval * time.Second)
		}
	}()
}
