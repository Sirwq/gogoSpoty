package botik

import (
	"sync"
	"time"
)

const cooldown = 10

type UserCooldowns struct {
	sync.Mutex
	m map[string]time.Time
}

func NewUserCooldowns() *UserCooldowns {
	return &UserCooldowns{
		m: map[string]time.Time{},
	}
}

func (uc *UserCooldowns) Store(key string, time time.Time) {
	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()
	uc.m[key] = time
}

func (uc *UserCooldowns) Load(key string) (time.Time, bool) {
	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	v, ok := uc.m[key]
	return v, ok
}
