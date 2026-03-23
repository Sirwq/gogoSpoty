//go:build standalone

package bot

import (
	"context"
	"gogoSpoty/internal/config"
	"sync"
)

type MemoryQueue struct {
	mu    sync.RWMutex
	queue []SongRequest
}

func NewQueue(cfg *config.RedisConfig) Queue {
	return &MemoryQueue{
		queue: make([]SongRequest, 0),
	}
}

func (q *MemoryQueue) Add(ctx context.Context, req SongRequest) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, req)
	return nil
}

func (q *MemoryQueue) Remove(context.Context) (SongRequest, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	var req SongRequest
	if len(q.queue) == 0 {
		return req, ErrQueueEmpty
	}
	req, q.queue = q.queue[0], q.queue[1:]
	return req, nil
}

func (q *MemoryQueue) Peek(ctx context.Context) (SongRequest, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	var req SongRequest
	if len(q.queue) == 0 {
		return req, ErrQueueEmpty
	}
	req = q.queue[0]

	return req, nil
}

func (q *MemoryQueue) Read(ctx context.Context) ([]SongRequest, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.queue, nil
}

func (q *MemoryQueue) Close() error {
	return nil
}
