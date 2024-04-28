package redis

import (
	cmap "github.com/orcaman/concurrent-map/v2"

	"eim/pkg/redis"
)

type Manager struct {
	rdsClient *redis.Client
	cache     cmap.ConcurrentMap[string, string]
}

func NewManager(endpoints []string, passwd string) (*Manager, error) {
	rdsClient, err := redis.NewClient(&redis.Config{
		Endpoints: endpoints,
		Password:  passwd,
	})
	if err != nil {
		return nil, err
	}
	return &Manager{rdsClient: rdsClient, cache: cmap.New[string]()}, nil
}
