package main

import (
	"context"
	"fmt"
	"time"
)

func r() {
	redisCTX := context.Background()

	client := NewRedisClient("1234")
	queue := NewQueue(client)

	req := SongRequest{
		Usename:     "vvxshin",
		TrackID:     "test123",
		TrackName:   "Nobody",
		TrackArtist: "Toxis",
		RequestedAt: time.Now(),
	}
	req2 := SongRequest{
		Usename:     "ASD",
		TrackID:     "test123",
		TrackName:   "SOmethingElse",
		TrackArtist: "Who",
		RequestedAt: time.Now(),
	}

	err := queue.Add(redisCTX, req)
	if err != nil {
		fmt.Println("Add error", err)
		return
	}

	err = queue.Add(redisCTX, req2)

	if err != nil {
		fmt.Println("Add error", err)
		return
	}

	songs, err := queue.Read(redisCTX)
	SongListPrinter(songs)

	s, err := queue.Remove(redisCTX)

	fmt.Println("Deleted song", s)

}
