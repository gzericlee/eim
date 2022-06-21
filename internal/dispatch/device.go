package dispatch

import (
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/storage"
	"eim/model"
)

type DeviceHandler struct {
	StorageRpc *storage.RpcClient
}

func (its *DeviceHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	global.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			device := &model.Device{}
			err := device.Deserialize(m.Body)
			if err != nil {
				global.Logger.Error("Error deserializing device", zap.Error(err))
				m.Finish()
				return
			}

			err = its.StorageRpc.SaveDevice(device)
			if err != nil {
				global.Logger.Error("Error saving device", zap.Error(err))
				m.Requeue(-1)
				return
			}

			m.Finish()
		}
	}(m))

	return nil
}
