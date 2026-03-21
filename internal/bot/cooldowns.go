package bot

import (
	"sync"
	"time"
)

const cooldown = 10

type UserCooldowns struct {
	sync.RWMutex
	m map[string]time.Time
}

func NewUserCooldowns() *UserCooldowns {
	return &UserCooldowns{
		m: map[string]time.Time{},
	}
}

func (uc *UserCooldowns) Store(key string, time time.Time) {
	uc.RWMutex.Lock()
	defer uc.RWMutex.Unlock()
	uc.m[key] = time
}

func (uc *UserCooldowns) Load(key string) (time.Time, bool) {
	uc.RWMutex.RLock()
	defer uc.RWMutex.RUnlock()

	v, ok := uc.m[key]
	return v, ok
}
