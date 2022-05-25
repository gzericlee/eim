package storage

import (
	"fmt"

	"github.com/nsqio/go-nsq"

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
				global.Logger.Warnf("Error inserting into Tidb: %v", err)
				m.Requeue(-1)
				return
			}

			err = redis.Set(fmt.Sprintf("%v:device:%v", device.UserId, device.DeviceId), m.Body)
			if err != nil {
				global.Logger.Warnf("Error saving into Redis cluster: %v", err)
				m.Requeue(-1)
				return
			}

			global.Logger.Infof("Store device: %v - %v", device.UserId, device.DeviceId)

			m.Finish()
		}
	}(m))

	return nil
}
