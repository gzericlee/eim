package redis

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

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

func TestManager_GetDevice(t *testing.T) {
	device, err := manager.GetDevice("user-1000", "device-1000")
	if err != nil {
		t.Error(err)
		return
	}
	body, _ := proto.Marshal(device)
	t.Log(string(body), err)
}

func TestManager_GetDevices(t *testing.T) {
	devices, err := manager.GetDevices("user-1")
	if err != nil {
		t.Error(err)
		return
	}
	body, _ := proto.Marshal(devices[0])
	t.Log(len(devices), string(body), err)
}

func TestManager_GetAllDevices(t *testing.T) {
	devices, err := manager.GetAllDevices()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(len(devices), err)
}

func TestManager_SaveUser(t *testing.T) {
	for i := 1; i <= 1000000; i++ {
		t.Log(manager.SaveBiz(&model.Biz{
			BizId:      fmt.Sprintf("user-%d", i),
			BizType:    model.Biz_USER,
			BizName:    fmt.Sprintf("用户-%d", i),
			TenantId:   "bingo",
			TenantName: "品高软件",
			Attributes: map[string]*anypb.Any{"password": {Value: []byte("pass@word1")}},
		}))
	}
}

func BenchmarkManager_SaveUser(b *testing.B) {
	b.N = 1000000
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = manager.SaveBiz(&model.Biz{
				BizId:      fmt.Sprintf("user-%d", b.N),
				BizName:    fmt.Sprintf("用户-%d", b.N),
				TenantId:   "bingo",
				TenantName: "品高软件",
				BizType:    model.Biz_USER,
				Attributes: map[string]*anypb.Any{"password": {Value: []byte("pass@word1")}},
			})
		}
	})
}

func TestManager_GetUser(t *testing.T) {
	t.Log(manager.GetBiz("user-1", "bingo"))
}

func TestManager_AppendBizMember(t *testing.T) {
	for i := 1; i <= 1000; i++ {
		t.Log(manager.AddBizMember(&model.BizMember{
			BizId:    "group-1",
			MemberId: fmt.Sprintf("user-%d", i),
			TenantId: "bingo",
		}))
	}
}

func TestManager_GetBizMembers(t *testing.T) {
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

func TestManager_SaveGateway(t *testing.T) {
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

func TestManager_GetGateways(t *testing.T) {
	gateways, err := manager.GetGateways()
	if err != nil {
		t.Error(err)
		return
	}
	for _, gateway := range gateways {
		t.Log(gateway.Ip, gateway.MemUsed, gateway.CpuUsed)
	}
}

func TestManager_SaveOfflineMessages(t *testing.T) {
	var msgs []*model.Message
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &model.Message{
			MsgId:    int64(i),
			FromId:   "user-1",
			ToId:     "user-2",
			SeqId:    1,
			Content:  "Hello, World!",
			SendTime: time.Now().Unix(),
		})
	}
	t.Log(manager.SaveOfflineMessages(msgs, "user-1", "device-1"))
	t.Log(manager.GetOfflineMessagesCount("user-1", "device-1"))
	t.Log(manager.GetOfflineMessages("user-1", "device-1"))
}

func BenchmarkManager_GetDevices(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idx := rand.Int64N(9999)
		_, _ = manager.GetDevices(fmt.Sprintf("user-%d", idx))
	}
}

func TestManager_RemoveDevice(t *testing.T) {
	keys, _ := manager.getAllKeys("devices:*")
	for _, key := range keys {
		manager.redisClient.Del(context.Background(), key).Err()
	}
}
