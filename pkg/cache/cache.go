package cache

import (
	"fmt"

	"github.com/dgraph-io/ristretto"
	"github.com/rcrowley/go-metrics"

	"eim/pkg/log"
)

type Cache struct {
	name  string
	cache *ristretto.Cache
}

func NewCache(name string, maxCost, numCounters int64, enableMetrics bool) (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		Metrics: enableMetrics,
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

	c := &Cache{
		name:  name,
		cache: cache,
	}

	if enableMetrics {
		metrics.Register(fmt.Sprintf("cache_%s_hits", name), metrics.NewFunctionalGauge(func() int64 {
			return int64(cache.Metrics.Hits())
		}))
		metrics.Register(fmt.Sprintf("cache_%s_misses", name), metrics.NewFunctionalGauge(func() int64 {
			return int64(cache.Metrics.Misses())
		}))
		metrics.Register(fmt.Sprintf("cache_%s_keys_added", name), metrics.NewFunctionalGauge(func() int64 {
			return int64(cache.Metrics.KeysAdded())
		}))
		metrics.Register(fmt.Sprintf("cache_%s_keys_evicted", name), metrics.NewFunctionalGauge(func() int64 {
			return int64(cache.Metrics.KeysEvicted())
		}))
		metrics.Register(fmt.Sprintf("cache_%s_cost_added", name), metrics.NewFunctionalGauge(func() int64 {
			return int64(cache.Metrics.CostAdded())
		}))
		metrics.Register(fmt.Sprintf("cache_%s_cost_evicted", name), metrics.NewFunctionalGauge(func() int64 {
			return int64(cache.Metrics.CostEvicted())
		}))
		metrics.Register(fmt.Sprintf("cache_%s_ratio", name), metrics.NewFunctionalGauge(func() int64 {
			return int64(cache.Metrics.Ratio())
		}))
	}

	return c, nil
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
