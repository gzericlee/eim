package redis

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"eim/internal/metric"
	"eim/internal/model"
)

var manager *Manager

func init() {
	var err error
	manager, err = NewManager(Config{
		RedisEndpoints:       []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
		RedisPassword:        "pass@word1",
		OfflineMessageExpire: time.Minute * 2,
		OfflineDeviceExpire:  time.Minute * 2,
	})
	if err != nil {
		panic(err)
	}
}

func TestSaveGroupMember(t *testing.T) {
	for i := 1; i <= 1000; i++ {
		t.Log(manager.AppendBizMember(&model.BizMember{
			BizId:   "group-1",
			BizType: "group",
			UserId:  "user-" + strconv.Itoa(i),
		}))
	}
}

func TestGetBizMembers(t *testing.T) {
	members, err := manager.GetBizMembers("group", "group-1")
	t.Log(len(members), members, err)
}

func BenchmarkGetBizMembers(b *testing.B) {
	b.N = 100000
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = manager.GetBizMembers("group", "group-1")
		}
	})
}

func TestSaveGateway(t *testing.T) {
	mMetric, err := metric.GetMachineMetric()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(manager.RegisterGateway(&model.Gateway{
		Ip:      "10.200.20.14",
		MemUsed: float32(mMetric.MemUsed),
		CpuUsed: float32(mMetric.CpuUsed),
	}, time.Second*10), err)
}

func TestGetGateways(t *testing.T) {
	gateways, err := manager.GetGateways()
	if err != nil {
		t.Error(err)
		return
	}
	for _, gateway := range gateways {
		t.Log(gateway.Ip, gateway.MemUsed, gateway.CpuUsed)
	}
}

func TestGetAll(t *testing.T) {
	all, err := manager.getAllValues(fmt.Sprintf("%s:device:*", "user-1000"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("count:", len(all))
}

func TestGetAllGateway(t *testing.T) {
	gateways, err := manager.GetGateways()
	if err != nil {
		t.Error(err)
		return
	}
	for _, gateway := range gateways {
		t.Log(gateway)
	}
	t.Log("Gateway节点数:", len(gateways))
}

func TestClearAll(t *testing.T) {
	err := manager.redisClient.(*redis.ClusterClient).ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
		_, err := client.FlushAll(ctx).Result()
		return err
	})
	if err != nil {
		t.Fatalf("清空redis失败: %v", err)
		return
	}
}
