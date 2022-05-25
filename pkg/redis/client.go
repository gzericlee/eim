package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	rdc *redis.ClusterClient
}

type Config struct {
	Endpoints []string
	Password  string
}

func NewClient(config *Config) (*Client, error) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          config.Endpoints,
		Password:       config.Password,
		ReadOnly:       true,
		RouteByLatency: true,
	})
	err := redisClient.Ping(context.Background()).Err()
	client := &Client{rdc: redisClient}
	return client, err
}

func (its *Client) Get(key string) (string, error) {
	return its.rdc.Get(context.Background(), key).Result()
}

func (its *Client) GetAll(key string) ([]string, error) {
	var result []string
	err := its.rdc.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
		keys, err := client.Keys(ctx, key).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			val, err := client.Get(context.Background(), key).Result()
			if err != nil {
				return err
			}
			result = append(result, val)
		}
		return nil
	})
	return result, err
}

func (its *Client) Set(key string, val interface{}, expiration time.Duration) error {
	return its.rdc.Set(context.Background(), key, val, expiration).Err()
}

func (its *Client) Incr(key string) (int64, error) {
	return its.rdc.Incr(context.Background(), key).Result()
}

func (its *Client) Publish(channel string, msg interface{}) {
	its.rdc.Publish(context.Background(), channel, msg)
}

func (its *Client) Subscribe(channel string, callback func(msg interface{}, err error)) {
	sub := its.rdc.Subscribe(context.Background(), channel)
	for {
		msg, err := sub.ReceiveMessage(context.Background())
		callback(msg, err)
	}
}
