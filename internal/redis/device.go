package redis

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/util/log"
)

func (its *Manager) SaveDevice(device *model.Device) error {
	body, err := proto.Marshal(device)
	if err != nil {
		return fmt.Errorf("proto marshal -> %w", err)
	}
	err = its.redisClient.Set(context.Background(), fmt.Sprintf("%v.device.%v", device.UserId, device.DeviceId), body, 0).Err()
	if err != nil {
		return fmt.Errorf("redis set -> %w", err)
	}
	return nil
}

func (its *Manager) GetUserDevices(userId string) ([]*model.Device, error) {
	values, err := its.getAll(fmt.Sprintf("%v.device.*", userId), 5000)
	if err != nil {
		return nil, fmt.Errorf("redis getAll -> %w", err)
	}
	var devices []*model.Device
	for _, value := range values {
		device := &model.Device{}
		err = proto.Unmarshal([]byte(value), device)
		if err != nil {
			log.Error("Error proto unmarshal. Drop it", zap.Error(err))
			continue
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (its *Manager) GetUserDevice(userId, deviceId string) (*model.Device, error) {
	value, err := its.redisClient.Get(context.Background(), fmt.Sprintf("%v.device.%v", userId, deviceId)).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get -> %w", err)
	}
	device := &model.Device{}
	err = proto.Unmarshal([]byte(value), device)
	if err != nil {
		return nil, fmt.Errorf("proto unmarshal -> %w", err)
	}
	return device, nil
}
