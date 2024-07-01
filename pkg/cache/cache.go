package cache

import (
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"

	"eim/util/log"
)

type Cache struct {
	cache *ristretto.Cache
}

func NewCache(name string, maxCost, numCounters int64) (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		Metrics: true,
		OnEvict: func(item *ristretto.Item) {
			log.Warn(fmt.Sprintf("cache evict %+v", item.Key))
		},
		OnReject: func(item *ristretto.Item) {
			log.Warn(fmt.Sprintf("cache reject %+v", item.Key))
		},
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("new ristretto cache -> %w", err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)
			log.Info(fmt.Sprintf("Cache `%s` metrics %+v", name, cache.Metrics))
		}
	}()

	return &Cache{
		cache: cache,
	}, nil
}

func (its *Cache) Set(key string, value interface{}) {
	ok := its.cache.Set(key, value, 1)
	if !ok {
		log.Warn("cache set failed")
	}
	its.cache.Wait()
}

func (its *Cache) Get(key string) (interface{}, bool) {
	return its.cache.Get(key)
}

func (its *Cache) Delete(key string) {
	its.cache.Del(key)
}

func (its *Cache) Close() {
	its.cache.Close()
}
