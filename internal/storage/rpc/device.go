package rpc

import (
	"context"
	"fmt"
	"time"

	"eim/internal/database"
	"eim/internal/model"
	"eim/internal/redis"
	"eim/pkg/cache"
	"eim/pkg/cache/notify"
	"eim/util/log"
)

type DeviceArgs struct {
	Device *model.Device
}

type DevicesReply struct {
	Devices []*model.Device
}

type DeviceReply struct {
	Device *model.Device
}

type Device struct {
	storageCache *cache.Cache
	redisManager *redis.Manager
	database     database.IDatabase
}

func (its *Device) SaveDevice(ctx context.Context, args *DeviceArgs, reply *DeviceReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.redisManager.SaveDevice(args.Device)
	if err != nil {
		return fmt.Errorf("save device -> %w", err)
	}

	key := fmt.Sprintf(cacheKeyFormat, deviceCachePool, args.Device.UserId, "*")
	err = notify.Del(deviceCachePool, key)
	if err != nil {
		return fmt.Errorf("del device(%s) cache -> %w", key, err)
	}

	return nil
}

func (its *Device) GetDevices(ctx context.Context, args *DeviceArgs, reply *DevicesReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	key := fmt.Sprintf(cacheKeyFormat, deviceCachePool, args.Device.UserId, "*")

	if cacheItem, exist := its.storageCache.Get(key); exist {
		reply.Devices = cacheItem.([]*model.Device)
		return nil
	}

	result, err, _ := group.Do(key, func() (interface{}, error) {
		devices, err := its.redisManager.GetDevices(args.Device.UserId)
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

func (its *Device) GetDevice(ctx context.Context, args *DeviceArgs, reply *DeviceReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	result := &DevicesReply{}

	err := its.GetDevices(ctx, args, result)
	if err != nil {
		return fmt.Errorf("get user devices -> %w", err)
	}

	for _, device := range result.Devices {
		if device.DeviceId == args.Device.DeviceId {
			reply.Device = device
			return nil
		}
	}

	return fmt.Errorf("device not found")
}
