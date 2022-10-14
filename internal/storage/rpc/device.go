package rpc

import (
	"context"

	"go.uber.org/zap"

	"eim/internal/redis"
	"eim/internal/types"
	"eim/pkg/log"
)

//type DeviceHandler struct{}
//
//func (its *DeviceHandler) HandleMessage(m *nsq.Message) error {
//	if len(m.Body) == 0 {
//		return nil
//	}
//
//	pool.SystemPool.Go(func(m *nsq.Message) func() {
//		return func() {
//			device := &types.Device{}
//			err := device.Deserialize(m.Body)
//			if err != nil {
//				m.Finish()
//				return
//			}
//
//			err = mainDb.SaveDevice(device)
//			if err != nil {
//				log.Error("Error inserting into Tidb", zap.Error(err))
//				m.Requeue(-1)
//				return
//			}
//
//			err = redis.Set(fmt.Sprintf("%v:device:%v", device.UserId, device.DeviceId), m.Body)
//			if err != nil {
//				log.Error("Error saving into Redis cluster", zap.Error(err))
//				m.Requeue(-1)
//				return
//			}
//
//			log.Info("Store device", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId))
//
//			m.Finish()
//		}
//	}(m))
//
//	return nil
//}

type DeviceRequest struct {
	Device *types.Device
}

type DeviceReply struct {
}

type Device struct {
}

func (its *Device) Save(ctx context.Context, req *DeviceRequest, reply *DeviceReply) error {
	err := mainDb.SaveDevice(req.Device)
	if err != nil {
		log.Error("Error inserting into Tidb", zap.Error(err))
		return err
	}

	err = redis.SaveUserDevice(req.Device)
	if err != nil {
		log.Error("Error saving into Redis cluster", zap.Error(err))
		return err
	}

	log.Info("Store device", zap.String("userId", req.Device.UserId), zap.String("deviceId", req.Device.DeviceId))

	return nil
}
