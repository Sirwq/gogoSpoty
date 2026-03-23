//go:build !standalone

package bot

import (
	"context"
	"encoding/json"
	"gogoSpoty/internal/config"
	"log"

	"github.com/redis/go-redis/v9"
)

const RedisKey = "song_queue"

type RedisQueue struct {
	client *redis.Client
}

func NewRedisClient(conf *config.RedisConfig) *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password,
			DB:       0,
		})

	return client
}

func NewQueue(cfg *config.RedisConfig) Queue {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
	})
	return &RedisQueue{client: client}
}

func (q *RedisQueue) Add(ctx context.Context, req SongRequest) error {
	data, err := json.Marshal(req)

	if err != nil {
		log.Println(err)
		return err
	}

	v := q.client.RPush(ctx, RedisKey, data)

	if v.Err() != nil {
		log.Println("Failed to push in REDIS queue")
		return v.Err()
	}

	return nil
}

func (q *RedisQueue) Remove(ctx context.Context) (SongRequest, error) {
	var req SongRequest

	v := q.client.LPop(ctx, RedisKey)

	if v.Err() == redis.Nil {
		return req, ErrQueueEmpty
	}

	if v.Err() != nil {
		log.Println("Failed to remove from REDIS queue")
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

func (q *RedisQueue) Read(ctx context.Context) ([]SongRequest, error) {
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

func (q *RedisQueue) Peek(ctx context.Context) (SongRequest, error) {
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

func (q *RedisQueue) Close() error {
	return q.client.Close()
}
