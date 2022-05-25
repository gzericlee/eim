package global

import (
	"github.com/patrickmn/go-cache"
)

var SystemCache *systemCache

type systemCache struct {
	cache *cache.Cache
}

func init() {
	SystemCache = &systemCache{}
	SystemCache.cache = cache.New(-1, -1)
}

func (its *systemCache) Save(key string, val interface{}) {
	its.cache.SetDefault(key, val)
}

func (its *systemCache) Get(key string) (interface{}, bool) {
	return its.cache.Get(key)
}
