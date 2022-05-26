package storage

import (
	"fmt"

	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/redis"
	"eim/model"
)

type DeviceHandler struct{}

func (its *DeviceHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	global.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			device := &model.Device{}
			err := device.Deserialize(m.Body)
			if err != nil {
				m.Finish()
				return
			}

			err = mainDb.SaveDevice(device)
			if err != nil {
				global.Logger.Error("Error inserting into Tidb", zap.Error(err))
				m.Requeue(-1)
				return
			}

			err = redis.Set(fmt.Sprintf("%v:device:%v", device.UserId, device.DeviceId), m.Body)
			if err != nil {
				global.Logger.Error("Error saving into Redis cluster", zap.Error(err))
				m.Requeue(-1)
				return
			}

			global.Logger.Info("Store device", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId))

			m.Finish()
		}
	}(m))

	return nil
}
