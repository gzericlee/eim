package oauth2

import (
	"sync"
	"time"
)

type Token map[string]interface{}
type TokenInfo map[string]interface{}
type UserInfo map[string]interface{}

type identity struct {
	UserInfo     UserInfo
	TokenInfo    TokenInfo
	AccessToken  string
	RefreshToken string
	ExpiresTime  time.Time
}

type identityMap struct {
	lock       *sync.RWMutex
	identities map[string]*identity
}

var identityManager *identityMap

func init() {
	identityManager = &identityMap{
		lock:       new(sync.RWMutex),
		identities: make(map[string]*identity),
	}
	go identityManager.AutoClean()
}

func (m *identityMap) Get(k string) *identity {
	m.lock.RLock()
	defer m.lock.RUnlock()
	auth := &identity{}
	if val, ok := m.identities[k]; ok {
		auth = val
	}
	return auth
}

func (m *identityMap) Set(k string, v *identity) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.identities[k] = v
}

func (m *identityMap) Delete(k string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exist := m.identities[k]; exist {
		delete(m.identities, k)
	}
}

func (m *identityMap) AutoClean() {
	for {
		time.Sleep(time.Minute * 5)
		m.lock.Lock()
		for k, v := range m.identities {
			if v.ExpiresTime.Before(time.Now()) {
				delete(m.identities, k)
			}
		}
		m.lock.Unlock()
	}
}
