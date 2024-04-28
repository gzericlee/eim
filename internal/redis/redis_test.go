package redis

import (
	"fmt"
	"strconv"
	"testing"

	"eim/internal/metric"
	"eim/internal/model"
)

var manager *Manager

func init() {
	var err error
	manager, err = NewManager([]string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"}, "pass@word1")
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
	body, _ := device.Serialize()
	t.Log(string(body), err)
}

func TestGetDevicesById(t *testing.T) {
	devices, err := manager.GetUserDevices("user-1")
	if err != nil {
		t.Error(err)
		return
	}
	body, _ := devices[0].Serialize()
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
		t.Log(manager.SaveGroupMember(&model.GroupMember{
			GroupId: "group-1",
			UserId:  "user-" + strconv.Itoa(i),
		}))
	}
}

func TestGetGroupMembers(t *testing.T) {
	members, err := manager.GetGroupMembers("group-1")
	t.Log(len(members), members, err)
}

func BenchmarkGetGroupMembers(b *testing.B) {
	b.N = 100000
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = manager.GetGroupMembers("group-1")
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
		MemUsed: mMetric.MemUsed,
		CpuUsed: mMetric.CpuUsed,
	}), err)
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
	users, err := manager.rdsClient.GetAll(fmt.Sprintf("%v*", "user-*:device:"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Redis总用户数:", len(users))
}

func TestGetAllGateway(t *testing.T) {
	gateways, err := manager.rdsClient.GetAll(fmt.Sprintf("%v:*", "gateway"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Gateway节点数:", len(gateways), gateways)
}
