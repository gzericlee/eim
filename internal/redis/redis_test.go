package redis

import (
	"testing"
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
