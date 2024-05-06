package rpc

import (
	"context"

	"go.uber.org/zap"

	"eim/internal/model"

	"eim/internal/database"
	"eim/internal/redis"
	"eim/pkg/log"
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
		log.Error("Error inserting into database", zap.Error(err))
		return err
	}

	err = its.RedisManager.SaveDevice(req.Device)
	if err != nil {
		log.Error("Error saving into redis cluster", zap.Error(err))
		return err
	}

	log.Debug("Store device", zap.String("userId", req.Device.UserId), zap.String("deviceId", req.Device.DeviceId), zap.String("gatewayIp", req.Device.GatewayIp), zap.Int32("state", req.Device.State))

	return nil
}
