package medis

import (
	"context"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Instance struct {
	mutex       sync.RWMutex
	lock        sync.Mutex
	redisClient *redis.ClusterClient
}

func NewInstance(redisClient *redis.ClusterClient) (*Instance, error) {
	instance := &Instance{redisClient: redisClient}
	current, err := instance.GetMax(MaxKey)
	if err != nil {
		return nil, err
	}
	if current == 0 {
		initMagazine(false)
		err = instance.SetMax(MaxKey, Capacity)
		if err != nil {
			return nil, err
		}
	} else {
		initMagazine(true)
	}
	return instance, nil
}

func (its *Instance) SetMax(key string, value int64) error {
	return its.redisClient.Set(context.Background(), key, value, 0).Err()
}

func (its *Instance) GetMax(key string) (int64, error) {
	result, err := its.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, err
}

func (its *Instance) Surplus(key string) (int64, error) {
	return its.redisClient.LLen(context.Background(), key).Result()
}

func (its *Instance) Push(supplement int64) error {
	its.lock.Lock()
	current, err := its.GetMax(MaxKey)
	if err != nil {
		return err
	}
	batch := int64(math.Floor(float64(supplement) / Unit))
	var i int64
	for i = 0; i < batch; i++ {
		for x := current + 1; x < (current + Unit + 1); x++ {
			sequence := generate(x)
			its.redisClient.LPush(context.Background(), ListKey, sequence)
		}
		its.redisClient.FlushAll(context.Background())
		current = current + Unit
	}
	err = its.SetMax(MaxKey, current)
	its.lock.Unlock()
	return err
}

func (its *Instance) RPop(channel chan int64, need int64) error {
	its.mutex.Lock()
	defer its.mutex.Unlock()

	ctx := context.Background()

	value, err := its.Surplus(ListKey)
	if err != nil {
		its.mutex.Unlock()
		return err
	}

	tx := its.redisClient.TxPipeline()
	tx.LRange(ctx, ListKey, value-need, value)
	tx.LTrim(ctx, ListKey, 0, value-need-1)
	cmds, err := tx.Exec(ctx)
	if err != nil {
		its.mutex.Unlock()
		return err
	}

	var sli []int64

	for _, cmd := range cmds {
		switch cmd.(type) {
		case *redis.StringSliceCmd:
			vals, _ := cmd.(*redis.StringSliceCmd).Result()
			for _, val := range vals {
				if val == "" {
					continue
				}
				mst, _ := strconv.ParseInt(val, 10, 64)
				sli = append(sli, mst)
			}
		default:
			continue
		}
	}
	sort.Slice(sli, func(i, j int) bool {
		return sli[i] < sli[j]
	})
	for _, mst := range sli {
		channel <- mst
	}
	return err
}

func (its *Instance) KvToChannel(channel chan int64, need, threshold int64) error {
	err := its.Supplement(threshold)
	if err != nil {
		return err
	}
	return its.RPop(channel, need)
}

func (its *Instance) Supplement(threshold int64) error {
	its.mutex.Lock()
	defer its.mutex.Unlock()

	value, err := its.Surplus(ListKey)
	if err != nil {
		return err
	}
	if value < threshold {
		err = its.Push(magazine.KvSupplement)
		if err != nil {
			return err
		}
	}
	return nil
}
