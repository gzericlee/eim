package redis

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
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
		Addrs:        config.Endpoints,
		Password:     config.Password,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		MaxRedirects: 8,
		PoolTimeout:  30 * time.Second,
	})

	err := redisClient.Ping(context.Background()).Err()
	client := &Client{rdc: redisClient}
	return client, err
}

func (its *Client) Get(key string) (string, error) {
	return its.rdc.Get(context.Background(), key).Result()
}

func (its *Client) GetAll(key string) ([]string, error) {
	var locker sync.RWMutex
	var result []string
	var cursor uint64
	err := its.rdc.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
		iter := client.Scan(ctx, cursor, key, 1000).Iterator()
		for iter.Next(ctx) {
			val, err := client.Get(context.Background(), iter.Val()).Result()
			if err != nil {
				return err
			}
			locker.Lock()
			result = append(result, val)
			locker.Unlock()
		}
		return iter.Err()
	})
	return result, err
}

func (its *Client) Set(key string, val interface{}, expiration time.Duration) error {
	return its.rdc.Set(context.Background(), key, val, expiration).Err()
}

func (its *Client) Incr(key string) (int64, error) {
	return its.rdc.Incr(context.Background(), key).Result()
}

func (its *Client) Decr(key string) (int64, error) {
	return its.rdc.Decr(context.Background(), key).Result()
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

func (its *Client) SAdd(key string, members ...interface{}) error {
	_, err := its.rdc.SAdd(context.Background(), key, members...).Result()
	return err
}

func (its *Client) SMembers(key string) ([]string, error) {
	return its.rdc.SMembers(context.Background(), key).Result()
}
