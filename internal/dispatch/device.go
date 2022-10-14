package dispatch

import (
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/internal/pool"
	storage_rpc "eim/internal/storage/rpc"
	"eim/internal/types"
	"eim/pkg/log"
)

type DeviceHandler struct {
	StorageRpc *storage_rpc.Client
}

func (its *DeviceHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	pool.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			device := &types.Device{}
			err := device.Deserialize(m.Body)
			if err != nil {
				log.Error("Error deserializing device", zap.Error(err))
				m.Finish()
				return
			}

			err = its.StorageRpc.SaveDevice(device)
			if err != nil {
				log.Error("Error saving device", zap.Error(err))
				m.Requeue(-1)
				return
			}

			m.Finish()
		}
	}(m))

	return nil
}
