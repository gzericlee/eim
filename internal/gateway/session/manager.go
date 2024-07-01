package session

import (
	"fmt"

	"eim/pkg/cache"
)

type Manager struct {
	prefix       string
	sessionCache *cache.Cache
}

func NewManager() *Manager {
	sessionCache, err := cache.NewCache("session", 3*1024*1024*1024, 1000000)
	if err != nil {
		panic(err)
	}
	return &Manager{
		prefix:       "session",
		sessionCache: sessionCache,
	}
}

func (its *Manager) Add(userId string, sess *Session) {
	key := fmt.Sprintf("%s:%s", its.prefix, userId)
	var sessions = map[string]*Session{}
	if value, exist := its.sessionCache.Get(key); exist {
		sessions = value.(map[string]*Session)
	}
	sessions[sess.device.DeviceId] = sess
	its.sessionCache.Set(key, sessions)
}

func (its *Manager) Get(userId string) map[string]*Session {
	key := fmt.Sprintf("%s:%s", its.prefix, userId)
	if value, exist := its.sessionCache.Get(key); exist {
		return value.(map[string]*Session)
	}
	return nil
}

func (its *Manager) Remove(userId, deviceId string) {
	key := fmt.Sprintf("%s:%s", its.prefix, userId)
	if value, exist := its.sessionCache.Get(key); exist {
		sessions := value.(map[string]*Session)
		delete(sessions, deviceId)
		its.sessionCache.Set(key, sessions)
	}
}
