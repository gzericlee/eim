package cache

import (
	"fmt"
	"time"

	"github.com/Yiling-J/theine-go"

	"github.com/gzericlee/eim/pkg/log"
)

type Cache[K comparable, V any] struct {
	name  string
	cache *theine.Cache[K, V]
}

func NewCache[K comparable, V any](name string, capacity int) (*Cache[K, V], error) {
	builder := theine.NewBuilder[K, V](int64(capacity))
	cache, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("new ristretto cache -> %w", err)
	}

	c := &Cache[K, V]{
		name:  name,
		cache: cache,
	}

	return c, nil
}

func (its *Cache[K, V]) Set(key K, value V) {
	ok := its.cache.Set(key, value, 1)
	if !ok {
		log.Warn("cache set failed")
	}
}

func (its *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	ok := its.cache.SetWithTTL(key, value, 1, ttl)
	if !ok {
		log.Warn("cache set failed")
	}
}

func (its *Cache[K, V]) Get(key K) (V, bool) {
	return its.cache.Get(key)
}

func (its *Cache[K, V]) Range(f func(key K, value V) bool) {
	its.cache.Range(f)
}

func (its *Cache[K, V]) Delete(key K) {
	its.cache.Delete(key)
}

func (its *Cache[K, V]) Close() {
	its.cache.Close()
}
