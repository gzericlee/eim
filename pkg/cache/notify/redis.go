package notify

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	ctx         context.Context
	redisClient redis.UniversalClient
}

func NewRedisProvider(client redis.UniversalClient) IProvider {
	return &RedisProvider{
		ctx:         context.Background(),
		redisClient: client,
	}
}

func (its *RedisProvider) OK() bool {
	_, err := its.redisClient.Ping(its.ctx).Result()
	return err == nil
}

func (its *RedisProvider) Pub(channel string, payload []string) error {
	_, err := its.redisClient.Publish(its.ctx, channel, strings.Join(payload, "#")).Result()
	if err != nil {
		return fmt.Errorf("redis publish(%s) -> %w", channel, err)
	}
	return nil
}

func (its *RedisProvider) Sub(channel string, callback func(payload []string)) error {
	msgChan := its.redisClient.Subscribe(its.ctx, channel).Channel()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			callback(strings.Split(msg.Payload, "#"))
		}
	}
}
