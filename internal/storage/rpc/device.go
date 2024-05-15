package rpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"eim/internal/model"
	"eim/util/log"

	"eim/internal/database"
	"eim/internal/redis"
)

type DeviceRequest struct {
	Device *model.Device
}

type DeviceReply struct {
}

type Device struct {
	Database     database.IDatabase
	RedisManager *redis.Manager
}

func (its *Device) Save(ctx context.Context, req *DeviceRequest, reply *DeviceReply) error {
	err := its.Database.SaveDevice(req.Device)
	if err != nil {
		return fmt.Errorf("save db device -> %w", err)
	}

	err = its.RedisManager.SaveDevice(req.Device)
	if err != nil {
		return fmt.Errorf("save redis device -> %w", err)
	}

	log.Debug("Store device", zap.String("userId", req.Device.UserId), zap.String("deviceId", req.Device.DeviceId), zap.String("gatewayIp", req.Device.GatewayIp), zap.Int32("state", req.Device.State))

	return nil
}
