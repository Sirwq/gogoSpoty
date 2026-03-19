package main

import (
	"context"
	"fmt"
)

func r() {
	redisCTX := context.Background()

	client := NewRedisClient("1234")
	queue := NewQueue(client)

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
