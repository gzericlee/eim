package redis

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/util/log"
)

const (
	userDeviceKeyFormat = "device:%s:%s"
	devicesKeyFormat    = "device:%s:*"
)

func (its *Manager) SaveDevice(device *model.Device) error {
	key := fmt.Sprintf(userDeviceKeyFormat, device.UserId, device.DeviceId)

	body, err := proto.Marshal(device)
	if err != nil {
		return fmt.Errorf("proto marshal -> %w", err)
	}

	err = its.redisClient.Set(context.Background(), key, body, 0).Err()
	if err != nil {
		return fmt.Errorf("redis set(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetDevices(userId string) ([]*model.Device, error) {
	key := fmt.Sprintf(devicesKeyFormat, userId)

	values, err := its.getAllValues(key)
	if err != nil {
		return nil, fmt.Errorf("redis getAllValues(%s) -> %w", key, err)
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

func (its *Manager) GetDevice(userId, deviceId string) (*model.Device, error) {
	key := fmt.Sprintf(userDeviceKeyFormat, userId, deviceId)

	value, err := its.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get(%s) -> %w", key, err)
	}

	device := &model.Device{}
	err = proto.Unmarshal([]byte(value), device)
	if err != nil {
		return nil, fmt.Errorf("proto unmarshal -> %w", err)
	}

	return device, nil
}
