package rpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"eim/internal/database"
	"eim/internal/model"
	"eim/internal/redis"
	"eim/pkg/cache"
	"eim/pkg/lock"
	"eim/pkg/log"
)

type DeviceArgs struct {
	Device *model.Device
}

type UserArgs struct {
	UserId   string
	TenantId string
	DeviceId string
}

type DevicesReply struct {
	Devices []*model.Device
}

type DeviceReply struct {
	Device *model.Device
}

type Device struct {
	lock         *lock.KeyLock
	storageCache *cache.Cache[string, []*model.Device]
	redisManager *redis.Manager
	database     database.IDatabase
}

func (its *Device) SaveDevice(ctx context.Context, args *DeviceArgs, reply *DeviceReply) error {
	err := its.redisManager.SaveDevice(args.Device)
	if err != nil {
		return fmt.Errorf("save device -> %w", err)
	}

	key := fmt.Sprintf(deviceCacheKeyFormat, deviceCachePool, args.Device.UserId, args.Device.TenantId, "*")

	err = storageRpc.RefreshDevicesCache(key, args.Device, ActionSave)
	if err != nil {
		log.Error("Error refresh device cache: %v", zap.Error(err))
	}

	return nil
}

func (its *Device) GetDevices(ctx context.Context, args *UserArgs, reply *DevicesReply) error {
	key := fmt.Sprintf(deviceCacheKeyFormat, deviceCachePool, args.UserId, args.TenantId, "*")

	_, unlock := its.lock.Lock(key, nil)
	defer unlock()

	if devices, exist := its.storageCache.Get(key); exist {
		reply.Devices = devices
		return nil
	}

	result, err, _ := singleGroup.Do(key, func() (interface{}, error) {
		devices, err := its.redisManager.GetDevices(args.UserId, args.TenantId)
		if err != nil {
			return nil, fmt.Errorf("get user devices -> %w", err)
		}
		return devices, nil
	})
	if err != nil {
		return fmt.Errorf("device single group do -> %w", err)
	}

	devices := result.([]*model.Device)
	reply.Devices = devices

	its.storageCache.Set(key, devices)

	return nil
}

func (its *Device) GetAllDevices(ctx context.Context, args *EmptyArgs, reply *DevicesReply) error {
	devices, err := its.redisManager.GetAllDevices()
	if err != nil {
		return fmt.Errorf("get all devices -> %w", err)
	}

	reply.Devices = devices

	return nil
}

func (its *Device) GetDevice(ctx context.Context, args *UserArgs, reply *DeviceReply) error {
	result := &DevicesReply{}

	err := its.GetDevices(ctx, args, result)
	if err != nil {
		return fmt.Errorf("get user devices -> %w", err)
	}

	for _, device := range result.Devices {
		if device.DeviceId == args.DeviceId {
			reply.Device = device
			return nil
		}
	}

	return fmt.Errorf("device not found")
}
