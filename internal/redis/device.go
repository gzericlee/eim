package redis

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/pkg/log"
)

const (
	devicesKeyFormat = "devices:%s:%s"
)

func (its *Manager) SaveDevice(device *model.Device) error {
	key := fmt.Sprintf(devicesKeyFormat, device.UserId, device.TenantId)

	body, err := proto.Marshal(device)
	if err != nil {
		return fmt.Errorf("proto marshal -> %w", err)
	}

	err = its.redisClient.HSet(context.Background(), key, device.DeviceId, body).Err()
	if err != nil {
		return fmt.Errorf("redis hset(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetDevices(userId, tenantId string) ([]*model.Device, error) {
	key := fmt.Sprintf(devicesKeyFormat, userId, tenantId)

	values, err := its.redisClient.HGetAll(context.Background(), key).Result()

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

func (its *Manager) GetAllDevices() ([]*model.Device, error) {
	key := fmt.Sprintf(devicesKeyFormat, "*", "*")

	values, err := its.redisClient.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis hgetall(%s) -> %w", key, err)
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

func (its *Manager) GetDevice(userId, tenantId, deviceId string) (*model.Device, error) {
	key := fmt.Sprintf(devicesKeyFormat, userId, tenantId)

	value, err := its.redisClient.HGet(context.Background(), key, deviceId).Result()
	if err != nil {
		return nil, fmt.Errorf("redis hget(%s) -> %w", key, err)
	}

	device := &model.Device{}
	err = proto.Unmarshal([]byte(value), device)
	if err != nil {
		return nil, fmt.Errorf("proto unmarshal -> %w", err)
	}

	return device, nil
}

func (its *Manager) RemoveDevice(userId, tenantId, deviceId string) error {
	key := fmt.Sprintf(devicesKeyFormat, userId, tenantId)

	err := its.redisClient.HDel(context.Background(), key, deviceId).Err()
	if err != nil {
		return fmt.Errorf("redis hdel(%s) -> %w", key, err)
	}

	return nil
}
