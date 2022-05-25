package redis

import (
	"eim/pkg/redis"
)

var rdsClient *redis.Client

func Set(key string, body []byte) error {
	err := rdsClient.Set(key, body, 0)
	return err
}

func InitRedisClusterClient(endpoints []string, passwd string) error {
	var err error
	rdsClient, err = redis.NewClient(&redis.Config{
		Endpoints: endpoints,
		Password:  passwd,
	})
	if err != nil {
		return err
	}
	return nil
}
