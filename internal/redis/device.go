package redis

import (
	"fmt"

	"eim/internal/types"
	"eim/pkg/json"
)

func SaveUserDevice(device *types.Device) error {
	body, _ := device.Serialize()
	return rdsClient.Set(fmt.Sprintf("%v:device:%v", device.UserId, device.DeviceId), body, 0)
}

func GetUserDevices(userId string) ([]*types.Device, error) {
	values, err := rdsClient.GetAll(fmt.Sprintf("%v:device:*", userId))
	if err != nil {
		return nil, err
	}
	var devices []*types.Device
	for _, value := range values {
		device := &types.Device{}
		err = device.Deserialize([]byte(value))
		if err != nil {
			continue
		}
		devices = append(devices, device)
	}
	return devices, err
}

func GetUserDevice(userId, deviceId string) (*types.Device, error) {
	value, err := rdsClient.Get(fmt.Sprintf("%v:device:%v", userId, deviceId))
	if err != nil {
		return nil, err
	}
	device := &types.Device{}
	err = json.Unmarshal([]byte(value), device)
	return device, err
}
