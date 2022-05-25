package seq

import (
	"sync"
	"testing"
	"time"

	"eim/internal/redis"
	"eim/model"
)

func init() {
	err := redis.InitRedisClusterClient([]string{"10.8.12.23:7001", "10.8.12.23:7002", "10.8.12.23:7003"}, "pass@word1")
	if err != nil {
		panic(err)
	}
}

func Test1(t *testing.T) {
	device := &model.Device{}
	device.DeviceId = "1"
	go func(device *model.Device) {
		time.Sleep(time.Second * 5)
		t.Log(device.DeviceId)
	}(device)
	device.DeviceId = "2"
	time.Sleep(time.Second * 10)
}

func TestID_Get(t *testing.T) {
	client, err := NewRpcClient("10.8.12.23:10000")
	if err != nil {
		panic(err)
	}
	wait := sync.WaitGroup{}
	for i := 0; i < 50000; i++ {
		wait.Add(1)
		go func() {
			t.Log(client.ID("user_1"))
			wait.Done()
		}()
	}
	wait.Wait()
}

func BenchmarkID_Get(b *testing.B) {
	client, err := NewRpcClient("10.8.12.23:10000")
	if err != nil {
		panic(err)
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client.ID("user_1")
		}
	})
}
