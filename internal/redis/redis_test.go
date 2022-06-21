package redis

import (
	"strconv"
	"testing"

	"eim/model"
)

func init() {
	_ = InitRedisClusterClient([]string{"10.8.12.23:7001", "10.8.12.23:7002", "10.8.12.23:7003", "10.8.12.23:7004", "10.8.12.23:7005"}, "pass@word1")
}

func TestGetDeviceByDeviceId(t *testing.T) {
	device, err := GetDeviceById("user_1000", "device_1000")
	body, _ := device.Serialize()
	t.Log(string(body), err)
}

func TestGetDevicesById(t *testing.T) {
	devices, err := GetDevicesById("user_1")
	if err != nil {
		panic(err)
	}
	body, _ := devices[0].Serialize()
	t.Log(len(devices), string(body), err)
}

func TestSaveUser(t *testing.T) {
	t.Log(SaveUser(&model.User{
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
		t.Log(SaveGroupMember(&model.GroupMember{
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
