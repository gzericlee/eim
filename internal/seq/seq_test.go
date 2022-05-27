package seq

import (
	"testing"

	"eim/internal/redis"
)

var rpcClient *RpcClient

func init() {
	err := redis.InitRedisClusterClient([]string{"10.8.12.23:7001", "10.8.12.23:7002", "10.8.12.23:7003"}, "pass@word1")
	if err != nil {
		panic(err)
	}

	rpcClient, err = NewRpcClient([]string{"10.8.12.23:2379", "10.8.12.23:2479", "10.8.12.23:2579"})
	if err != nil {
		panic(err)
	}
}

func TestID_Get(t *testing.T) {
	s := 0
	for i := 0; i < 150000; i++ {
		_, err := rpcClient.Id("user_1")
		if err != nil {
			t.Log(err)
		} else {
			s++
		}
	}
	t.Log(s)
}

func BenchmarkID_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rpcClient.Id("user_1")
		}
	})
}
