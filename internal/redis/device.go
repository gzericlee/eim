package redis

import (
	"fmt"

	"eim/internal/model"
	"eim/pkg/json"
)

func (its *Manager) SaveDevice(device *model.Device) error {
	body, _ := device.Serialize()
	return its.rdsClient.Set(fmt.Sprintf("%v:device:%v", device.UserId, device.DeviceId), body, 0)
}

func (its *Manager) GetUserDevices(userId string) ([]*model.Device, error) {
	values, err := its.rdsClient.GetAll(fmt.Sprintf("%v:device:*", userId))
	if err != nil {
		return nil, err
	}
	var devices []*model.Device
	for _, value := range values {
		device := &model.Device{}
		err = device.Deserialize([]byte(value))
		if err != nil {
			continue
		}
		devices = append(devices, device)
	}
	return devices, err
}

func (its *Manager) GetUserDevice(userId, deviceId string) (*model.Device, error) {
	value, err := its.rdsClient.Get(fmt.Sprintf("%v:device:%v", userId, deviceId))
	if err != nil {
		return nil, err
	}
	device := &model.Device{}
	err = json.Unmarshal([]byte(value), device)
	return device, err
}
