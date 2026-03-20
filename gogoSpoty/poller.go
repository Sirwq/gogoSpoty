package main

import (
	"context"
	"gogoSpoty/botik"
	"gogoSpoty/spoty"
	"log"
	"time"

	"github.com/zmb3/spotify/v2"
)

type Poller struct {
	Client   *spotify.Client
	Track    *spoty.Track
	Queue    *botik.Queue
	Interval time.Duration
}

func (p *Poller) Start(ctx context.Context) {
	go func() {
		for {
			playing, err := p.Client.PlayerCurrentlyPlaying(ctx)

			if err != nil {
				log.Printf("Error fetching track %v\n", err)
				time.Sleep(p.Interval * time.Second)
				continue
			}

			SongRequest, err := p.Queue.Peek(ctx)

			if err == botik.ErrQueueEmpty {

			} else if err != nil {
				log.Printf("Error in poller: %v\n", err)
			}

			if playing != nil && playing.Item != nil {
				spoty.UpdateTrack(p.Track, playing)
				if playing.Item.ID == spotify.ID(SongRequest.TrackID) {
					_, err := p.Queue.Remove(ctx)
					if err != nil {
						log.Printf("MOZG UMER ERRROR CHANGE LATER")
					}
				} else {
					// song, err := p.Queue.Remove(ctx)

					// err = spotify.Client.QueueS
				}
			}

			time.Sleep(p.Interval * time.Second)
		}
	}()
}
