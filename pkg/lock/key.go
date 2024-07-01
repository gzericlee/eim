package lock

import "sync"

type KeyLock struct {
	locks map[string]*sync.Mutex
	mutex sync.Mutex
}

func NewKeyLock() *KeyLock {
	return &KeyLock{
		locks: make(map[string]*sync.Mutex),
	}
}

func (kl *KeyLock) Lock(key string) {
	kl.mutex.Lock()
	lock, ok := kl.locks[key]
	if !ok {
		lock = &sync.Mutex{}
		kl.locks[key] = lock
	}
	kl.mutex.Unlock()

	lock.Lock()
}

func (kl *KeyLock) TryLock(key string) bool {
	kl.mutex.Lock()
	lock, ok := kl.locks[key]
	if !ok {
		lock = &sync.Mutex{}
		kl.locks[key] = lock
	}
	kl.mutex.Unlock()

	return lock.TryLock()
}

func (kl *KeyLock) Unlock(key string) {
	kl.mutex.Lock()
	lock, ok := kl.locks[key]
	if ok {
		delete(kl.locks, key)
		lock.Unlock()
	}
	kl.mutex.Unlock()
}
