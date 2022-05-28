package redis

import (
	"testing"

	"github.com/google/uuid"

	"eim/model"
)

func init() {
	_ = InitRedisClusterClient([]string{"10.8.12.23:7001", "10.8.12.23:7002", "10.8.12.23:7003"}, "pass@word1")
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
		UserId:   uuid.New().String(),
		LoginId:  "lirui",
		UserName: "李锐",
		Password: "pass@word1",
		Company:  "bingo",
		SeqId:    100,
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
