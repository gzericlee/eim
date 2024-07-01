package lock

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	mutex *redsync.Mutex
}

func NewRedisLock(client redis.UniversalClient, name string) *RedisLock {
	return &RedisLock{
		mutex: redsync.New(goredis.NewPool(client)).NewMutex(name),
	}
}

func (its *RedisLock) Lock() error {
	return its.mutex.Lock()
}

func (its *RedisLock) Unlock() error {
	_, err := its.mutex.Unlock()
	return err
}

func (its *RedisLock) TryLock() error {
	return its.mutex.TryLock()
}
