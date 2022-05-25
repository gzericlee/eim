package redis

import (
	"fmt"

	"eim/model"
	"eim/pkg/json"
)

func GetDevicesById(userId string) ([]*model.Device, error) {
	values, err := rdsClient.GetAll(fmt.Sprintf("%v:device:*", userId))
	if err != nil {
		return nil, err
	}
	var devices []*model.Device
	for _, value := range values {
		device := &model.Device{}
		err = json.Unmarshal([]byte(value), device)
		if err != nil {
			continue
		}
		devices = append(devices, device)
	}
	return devices, err
}

func GetDeviceById(userId, deviceId string) (*model.Device, error) {
	value, err := rdsClient.Get(fmt.Sprintf("%v:device:%v", userId, deviceId))
	if err != nil {
		return nil, err
	}
	device := &model.Device{}
	err = json.Unmarshal([]byte(value), device)
	return device, err
}
