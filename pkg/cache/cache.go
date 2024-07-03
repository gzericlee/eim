package cache

import (
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/jedib0t/go-pretty/v6/table"

	"eim/util/log"
)

type Cache struct {
	name  string
	cache *ristretto.Cache
}

func NewCache(name string, maxCost, numCounters int64, metrics bool) (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		Metrics: metrics,
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

	if metrics {
		go c.printMetrics()
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

func (its *Cache) printMetrics() {
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			{
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Cache Name", "Metrics"})
				t.AppendRows([]table.Row{{
					its.name,
					its.cache.Metrics},
				})
				t.AppendSeparator()
				t.Render()
			}
		}
	}
}
