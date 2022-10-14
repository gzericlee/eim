package redis

import (
	"strconv"
	"testing"

	"eim/internal/metric"
	"eim/internal/types"
)

func init() {
	_ = InitRedisClusterClient([]string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"}, "pass@word1")
}

func TestGetDeviceByDeviceId(t *testing.T) {
	device, err := GetUserDevice("user_1000", "device_1000")
	body, _ := device.Serialize()
	t.Log(string(body), err)
}

func TestGetDevicesById(t *testing.T) {
	devices, err := GetUserDevices("user_1")
	if err != nil {
		panic(err)
	}
	body, _ := devices[0].Serialize()
	t.Log(len(devices), string(body), err)
}

func TestSaveUser(t *testing.T) {
	t.Log(SaveUser(&types.User{
		UserId:   "1",
		LoginId:  "lirui",
		UserName: "李锐",
		Password: "pass@word1",
		Company:  "bingo",
	}))
}

func TestGetUser(t *testing.T) {
	user, err := GetUser("lirui", "bingo")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user.UserId, user.LoginId, user.UserName)
}

func TestSaveGroupMember(t *testing.T) {
	for i := 1; i <= 1000; i++ {
		t.Log(SaveGroupMember(&types.GroupMember{
			GroupId: "group-1",
			UserId:  "user-" + strconv.Itoa(i),
		}))
	}
}

func TestGetGroupMembers(t *testing.T) {
	members, err := GetGroupMembers("group-1")
	t.Log(len(members), members, err)
}

func BenchmarkGetGroupMembers(b *testing.B) {
	b.N = 100000
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GetGroupMembers("group-1")
		}
	})
}

func TestSaveGateway(t *testing.T) {
	mMetric, err := metric.GetMachineMetric()
	t.Log(RegisterGateway(&types.Gateway{
		Ip:      "10.200.20.14",
		MemUsed: mMetric.MemUsed,
		CpuUsed: mMetric.CpuUsed,
	}), err)
}

func TestGetGateways(t *testing.T) {
	gateways, err := GetGateways()
	if err != nil {
		t.Error(err)
		return
	}
	for _, gateway := range gateways {
		t.Log(gateway.Ip, gateway.MemUsed, gateway.CpuUsed)
	}
}
