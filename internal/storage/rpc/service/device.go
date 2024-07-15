package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/gzericlee/eim/internal/database"
	"github.com/gzericlee/eim/internal/model"
	rpcclient "github.com/gzericlee/eim/internal/storage/rpc/client"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
	"github.com/gzericlee/eim/pkg/cache"
	"github.com/gzericlee/eim/pkg/lock"
	"github.com/gzericlee/eim/pkg/log"
)

type DeviceService struct {
	lock         *lock.KeyLock
	storageCache *cache.Cache[string, []*model.Device]
	database     database.IDatabase
	refreshRpc   *rpcclient.RefresherClient
}

func NewDeviceService(lock *lock.KeyLock, storageCache *cache.Cache[string, []*model.Device], database database.IDatabase, refreshRpc *rpcclient.RefresherClient) *DeviceService {
	return &DeviceService{
		lock:         lock,
		storageCache: storageCache,
		database:     database,
		refreshRpc:   refreshRpc,
	}
}

func (its *DeviceService) SaveDevice(ctx context.Context, args *rpcmodel.DeviceArgs, reply *rpcmodel.DeviceReply) error {
	isNew, err := its.database.SaveDevice(args.Device)
	if err != nil {
		return fmt.Errorf("save device -> %w", err)
	}
	action := rpcmodel.ActionUpdate
	if isNew {
		action = rpcmodel.ActionInsert
	}

	key := fmt.Sprintf(deviceCacheKeyFormat, deviceCachePool, args.Device.UserId, args.Device.TenantId, "*")

	err = its.refreshRpc.RefreshDevicesCache(key, args.Device, action)
	if err != nil {
		log.Error("Error refresh device cache: %v", zap.Error(err))
	}

	return nil
}

func (its *DeviceService) InsertDevice(ctx context.Context, args *rpcmodel.DeviceArgs, reply *rpcmodel.DeviceReply) error {
	err := its.database.InsertDevice(args.Device)
	if err != nil {
		return fmt.Errorf("insert device -> %w", err)
	}

	key := fmt.Sprintf(deviceCacheKeyFormat, deviceCachePool, args.Device.UserId, args.Device.TenantId, "*")

	err = its.refreshRpc.RefreshDevicesCache(key, args.Device, rpcmodel.ActionInsert)
	if err != nil {
		log.Error("Error refresh device cache: %v", zap.Error(err))
	}

	return nil
}

func (its *DeviceService) UpdateDevice(ctx context.Context, args *rpcmodel.DeviceArgs, reply *rpcmodel.DeviceReply) error {
	err := its.database.UpdateDevice(args.Device)
	if err != nil {
		return fmt.Errorf("update device -> %w", err)
	}

	key := fmt.Sprintf(deviceCacheKeyFormat, deviceCachePool, args.Device.UserId, args.Device.TenantId, "*")

	err = its.refreshRpc.RefreshDevicesCache(key, args.Device, rpcmodel.ActionUpdate)
	if err != nil {
		log.Error("Error refresh device cache: %v", zap.Error(err))
	}

	return nil
}

func (its *DeviceService) GetDevices(ctx context.Context, args *rpcmodel.UserArgs, reply *rpcmodel.DevicesReply) error {
	key := fmt.Sprintf(deviceCacheKeyFormat, deviceCachePool, args.UserId, args.TenantId, "*")

	_, unlock := its.lock.Lock(key, nil)
	defer unlock()

	if devices, exist := its.storageCache.Get(key); exist {
		reply.Devices = devices
		return nil
	}

	result, err, _ := singleGroup.Do(key, func() (interface{}, error) {
		devices, err := its.database.GetDevicesByUser(args.UserId, args.TenantId)
		if err != nil {
			return nil, fmt.Errorf("get user devices -> %w", err)
		}
		return devices, nil
	})
	if err != nil {
		return fmt.Errorf("single group do -> %w", err)
	}

	devices := result.([]*model.Device)
	reply.Devices = devices

	its.storageCache.Set(key, devices)

	return nil
}

func (its *DeviceService) GetDevice(ctx context.Context, args *rpcmodel.UserArgs, reply *rpcmodel.DeviceReply) error {
	result := &rpcmodel.DevicesReply{}

	err := its.GetDevices(ctx, args, result)
	if err != nil {
		return fmt.Errorf("get devices -> %w", err)
	}

	for _, device := range result.Devices {
		if device.DeviceId == args.DeviceId {
			reply.Device = device
		}
	}

	return nil
}
