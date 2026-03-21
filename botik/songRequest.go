package botik

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrQueueEmpty = errors.New("queue is empty")

const RedisKey = "song_queue"

type SongRequest struct {
	Usename     string    `json:"username"`
	TrackID     string    `json:"track_id"`
	TrackName   string    `json:"track_name"`
	TrackArtist string    `json:"track_artist"`
	RequestedAt time.Time `json:"requested_at"`
}

type Queue struct {
	client *redis.Client
}

func (sr SongRequest) String() string {
	return fmt.Sprintf("username: %s\nTrack: %s - %s\nRequested at: %v", sr.Usename, sr.TrackName, sr.TrackArtist, sr.RequestedAt)
}

func NewRedisClient(pass string) *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: pass,
			DB:       0,
		})

	return client
}

func NewQueue(client *redis.Client) *Queue {
	return &Queue{client: client}
}

func (q *Queue) Add(ctx context.Context, req SongRequest) error {
	data, err := json.Marshal(req)

	if err != nil {
		fmt.Printf("Error marshaling SONG REQUEST for REDIS")
		return err
	}

	v := q.client.RPush(ctx, RedisKey, data)

	if v.Err() != nil {
		fmt.Println("Error while pushing to REDIS QUEUE")
		return v.Err()
	}

	return nil
}

func (q *Queue) Remove(ctx context.Context) (SongRequest, error) {
	var req SongRequest

	v := q.client.LPop(ctx, RedisKey)

	if v.Err() == redis.Nil {
		return req, nil
	}

	if v.Err() != nil {
		fmt.Println("Error while removing from redis queue")
		return req, v.Err()
	}

	data, err := v.Result()

	if err != nil {
		return req, err
	}

	err = json.Unmarshal([]byte(data), &req)

	if err != nil {
		return req, err
	}
	return req, nil
}

func (q *Queue) Read(ctx context.Context) ([]SongRequest, error) {
	songs, err := q.client.LRange(ctx, RedisKey, 0, -1).Result()

	if err != nil {
		return nil, err
	}

	result := make([]SongRequest, 0, len(songs))

	for _, song := range songs {
		var req SongRequest
		if err := json.Unmarshal([]byte(song), &req); err != nil {
			return nil, err
		}
		result = append(result, req)
	}
	return result, nil
}

func SongListPrinter(srlist []SongRequest) {
	for _, song := range srlist {
		fmt.Println(song)
	}
}

func (q *Queue) Peek(ctx context.Context) (SongRequest, error) {
	var s SongRequest

	data, err := q.client.LIndex(ctx, RedisKey, 0).Result()

	if err == redis.Nil {
		return s, ErrQueueEmpty
	}

	if err != nil {
		return s, err
	}

	err = json.Unmarshal([]byte(data), &s)
	if err != nil {
		return s, err
	}

	return s, nil

}
