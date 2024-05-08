package redis

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	"eim/internal/model"
)

func (its *Manager) SaveDevice(device *model.Device) error {
	body, err := proto.Marshal(device)
	if err != nil {
		return err
	}
	return its.redisClient.Set(context.Background(), fmt.Sprintf("%v.device.%v", device.UserId, device.DeviceId), body, 0).Err()
}

func (its *Manager) GetUserDevices(userId string) ([]*model.Device, error) {
	values, err := its.getAll(fmt.Sprintf("%v.device.*", userId), 5000)
	if err != nil {
		return nil, err
	}
	var devices []*model.Device
	for _, value := range values {
		device := &model.Device{}
		err = proto.Unmarshal([]byte(value), device)
		if err != nil {
			continue
		}
		devices = append(devices, device)
	}
	return devices, err
}

func (its *Manager) GetUserDevice(userId, deviceId string) (*model.Device, error) {
	value, err := its.redisClient.Get(context.Background(), fmt.Sprintf("%v.device.%v", userId, deviceId)).Result()
	if err != nil {
		return nil, err
	}
	device := &model.Device{}
	err = proto.Unmarshal([]byte(value), device)
	return device, err
}
