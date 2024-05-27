package notify

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisNotifier struct {
	ctx      context.Context
	redisCli redis.UniversalClient
}

func NewRedisNotifier(r redis.UniversalClient) Notifier {
	return &RedisNotifier{
		ctx:      context.Background(),
		redisCli: r,
	}
}

func (g *RedisNotifier) OK() bool {
	_, err := g.redisCli.Ping(g.ctx).Result()
	return err == nil
}

func (g *RedisNotifier) Pub(channel, payload string) error {
	_, err := g.redisCli.Publish(g.ctx, channel, payload).Result()
	if err != nil {
		return fmt.Errorf("redis publish(%s) -> %w", channel, err)
	}
	return nil
}

func (g *RedisNotifier) Sub(channel string, callback func(payload string)) error {
	msgChan := g.redisCli.Subscribe(g.ctx, channel).Channel()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			callback(msg.Payload)
		}
	}
}
