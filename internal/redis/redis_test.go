package redis

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"

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

func TestGetDeviceByDeviceId(t *testing.T) {
	device, err := manager.GetUserDevice("user-1000", "device-1000")
	if err != nil {
		t.Error(err)
		return
	}
	body, _ := proto.Marshal(device)
	t.Log(string(body), err)
}

func TestGetDevicesById(t *testing.T) {
	devices, err := manager.GetUserDevices("user-1")
	if err != nil {
		t.Error(err)
		return
	}
	body, _ := proto.Marshal(devices[0])
	t.Log(len(devices), string(body), err)
}

func TestSaveUser(t *testing.T) {
	t.Log(manager.SaveUser(&model.User{
		UserId:     "1",
		LoginId:    "lirui",
		UserName:   "李锐",
		Password:   "pass@word1",
		TenantId:   "bingo",
		TenantName: "品高软件",
	}))
}

func TestGetUser(t *testing.T) {
	user, err := manager.GetUser("lirui", "bingo")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user.UserId, user.LoginId, user.UserName)
}

func TestSaveGroupMember(t *testing.T) {
	for i := 1; i <= 1000; i++ {
		t.Log(manager.SaveBizMember(&model.BizMember{
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
	all, err := manager.getAllValues(fmt.Sprintf("%s*", "user-"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("count:", len(all))
}

func TestGetAllGateway(t *testing.T) {
	gateways, err := manager.getAllValues(fmt.Sprintf("gateway.*"))
	if err != nil {
		t.Error(err)
		return
	}
	for _, gateway := range gateways {
		t.Log(gateway)
	}
	t.Log("Gateway节点数:", len(gateways))
}

func TestManager_GetOfflineMessageCount(t *testing.T) {
	count, err := manager.GetOfflineMessageCount("user-1", "device-1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(count)
}

func TestManager_GetOfflineMessagesByDevice(t *testing.T) {
	messages, err := manager.GetOfflineMessagesByDevice("user-1", "device-1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(messages)
}
