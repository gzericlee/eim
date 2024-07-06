package session

import (
	"fmt"

	"eim/pkg/cache"
)

type Manager struct {
	prefix       string
	sessionCache *cache.Cache[string, map[string]*Session]
}

func NewManager() *Manager {
	sessionCache, err := cache.NewCache[string, map[string]*Session]("sessions", 1000000)
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
		sessions = value
	}
	sessions[sess.device.DeviceId] = sess
	its.sessionCache.Set(key, sessions)
}

func (its *Manager) Get(userId string) map[string]*Session {
	key := fmt.Sprintf("%s:%s", its.prefix, userId)
	if sessions, exist := its.sessionCache.Get(key); exist {
		return sessions
	}
	return nil
}

func (its *Manager) Remove(userId, deviceId string) {
	key := fmt.Sprintf("%s:%s", its.prefix, userId)
	if sessions, exist := its.sessionCache.Get(key); exist {
		delete(sessions, deviceId)
		its.sessionCache.Set(key, sessions)
	}
}
