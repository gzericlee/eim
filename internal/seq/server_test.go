package seq

import (
	"log"
	"sync"
	"testing"
	"time"

	"eim/internal/seq/rpc"
)

var rpcClient *rpc.Client
var etcdEndpoints = []string{"127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"}

func init() {
	go func() {
		err := rpc.StartServer(rpc.Config{
			Ip:             "127.0.0.1",
			Port:           18080,
			EtcdEndpoints:  etcdEndpoints,
			RedisEndpoints: []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
			RedisPassword:  "pass@word1",
		})
		log.Println(err)
	}()

	time.Sleep(time.Second * 5)

	var err error
	rpcClient, err = rpc.NewClient(etcdEndpoints)
	if err != nil {
		panic(err)
	}
}

func TestID_Get(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			seq, err := rpcClient.Number("user_2")
			if err != nil {
				t.Log(err)
			}
			t.Log(seq)
		}()
	}
	wg.Wait()
}

func BenchmarkID_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := rpcClient.Number("user_2")
			if err != nil {
				b.Error(err)
				return
			}
		}
	})
}
