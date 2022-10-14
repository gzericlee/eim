package seq

import (
	"testing"

	"eim/internal/redis"
	"eim/internal/seq/rpc"
)

var rpcClient *rpc.Client

func init() {
	err := redis.InitRedisClusterClient([]string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003"}, "pass@word1")
	if err != nil {
		panic(err)
	}

	rpcClient, err = rpc.NewClient([]string{"127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"})
	if err != nil {
		panic(err)
	}
}

func TestID_Get(t *testing.T) {
	s := 0
	for i := 0; i < 150000; i++ {
		_, err := rpcClient.Number("user_1")
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
			_, _ = rpcClient.Number("user_1")
		}
	})
}
