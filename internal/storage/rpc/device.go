package rpc

import (
	"context"
	"fmt"
	"time"

	"eim/internal/database"
	"eim/internal/model"
	"eim/pkg/cache"
	"eim/pkg/cache/notify"
	"eim/util/log"
)

type SaveDeviceArgs struct {
	Device *model.Device
}

type GetDevicesArgs struct {
	UserId string
}

type GetDeviceArgs struct {
	UserId   string
	DeviceId string
}

type DevicesReply struct {
	Devices []*model.Device
}

type DeviceReply struct {
	Device *model.Device
}

type Device struct {
	storageCache *cache.Cache
	database     database.IDatabase
}

func (its *Device) SaveDevice(ctx context.Context, args *SaveDeviceArgs, reply *DeviceReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.database.SaveDevice(args.Device)
	if err != nil {
		return fmt.Errorf("save device -> %w", err)
	}

	key := fmt.Sprintf("%s:%s", deviceCachePool, args.Device.DeviceId)
	err = notify.Del(deviceCachePool, key)
	if err != nil {
		return fmt.Errorf("del device(%s) cache -> %w", key, err)
	}

	return nil
}

func (its *Device) GetDevices(ctx context.Context, args *GetDevicesArgs, reply *DevicesReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	key := fmt.Sprintf("%s:%s", deviceCachePool, args.UserId)

	if cacheItem, exist := its.storageCache.Get(key); exist {
		reply.Devices = cacheItem.([]*model.Device)
		return nil
	}

	result, err, _ := group.Do(key, func() (interface{}, error) {
		devices, err := its.database.GetDevices(args.UserId)
		if err != nil {
			return nil, fmt.Errorf("get user devices -> %w", err)
		}
		return devices, nil
	})
	if err != nil {
		return fmt.Errorf("group do -> %w", err)
	}

	devices := result.([]*model.Device)
	its.storageCache.Put(key, devices)

	reply.Devices = devices

	return nil
}

func (its *Device) GetDevice(ctx context.Context, args *GetDeviceArgs, reply *DeviceReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	result := &DevicesReply{}

	err := its.GetDevices(ctx, &GetDevicesArgs{UserId: args.UserId}, result)
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
