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
	Client       *spotify.Client
	Track        *spoty.Track
	Queue        *botik.Queue
	Interval     time.Duration
	LastQueuedID string
}

func NewPoller(client *spotify.Client, track *spoty.Track, q *botik.Queue, interval time.Duration) *Poller {
	return &Poller{
		Client:   client,
		Track:    track,
		Queue:    q,
		Interval: interval,
	}
}

func (p *Poller) Start(ctx context.Context) {
	for {
		playing, err := p.Client.PlayerCurrentlyPlaying(ctx)

		if err != nil {
			log.Printf("Error fetching track %v\n", err)
			p.wait()
			continue
		}

		if playing == nil || playing.Item == nil {
			p.wait()
			continue
		}

		spoty.UpdateTrack(p.Track, playing)
		SongRequest, err := p.Queue.Peek(ctx)

		if err == botik.ErrQueueEmpty {
			p.wait()
			continue
		} else if err != nil {
			log.Printf("Error in poller: %v\n", err)
			p.wait()
			continue
		}

		if playing.Item.ID == spotify.ID(SongRequest.TrackID) {
			_, err := p.Queue.Remove(ctx)
			logErr(err, "Error removing from queue in poller")
			p.wait()
			continue
		}

		if p.LastQueuedID == SongRequest.TrackID {
			p.wait()
			continue
		}

		err = p.Client.QueueSong(ctx, spotify.ID(SongRequest.TrackID))
		p.LastQueuedID = SongRequest.TrackID
		logErr(err, "Error queuing song in poller")
		p.wait()
	}
}

func logErr(err error, msg string) {
	if err != nil {
		log.Println(msg, "Error:"+err.Error())
	}
}

func (p *Poller) wait() {
	time.Sleep(p.Interval * time.Second)
}
