package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type DeviceClient struct {
	*rpcxclient.XClientPool
}

func (its *DeviceClient) SaveDevice(device *model.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "SaveDevice", &rpcmodel.DeviceArgs{Device: device}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call SaveDevice -> %w", err)
	}

	return nil
}

func (its *DeviceClient) InsertDevice(device *model.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "InsertDevice", &rpcmodel.DeviceArgs{Device: device}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertDevice -> %w", err)
	}

	return nil
}

func (its *DeviceClient) UpdateDevice(device *model.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "UpdateDevice", &rpcmodel.DeviceArgs{Device: device}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call UpdateDevice -> %w", err)
	}

	return nil
}

func (its *DeviceClient) GetDevice(userId, tenantId, deviceId string) (*model.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.DeviceReply{}
	err := its.Get().Call(ctx, "GetDevice", &rpcmodel.UserArgs{UserId: userId, TenantId: tenantId, DeviceId: deviceId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetDevice -> %w", err)
	}

	return reply.Device, nil
}

func (its *DeviceClient) GetDevices(userId, tenantId string) ([]*model.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.DevicesReply{}
	err := its.Get().Call(ctx, "GetDevices", &rpcmodel.UserArgs{UserId: userId, TenantId: tenantId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetDevices -> %w", err)
	}

	return reply.Devices, nil
}
