package dispatch

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"eim/internal/model"
	storagerpc "eim/internal/storage/rpc"
	"eim/util/log"
)

const maxRetries = 3
const batchSize = 500

type saveTask struct {
	messages   chan *nats.Msg
	storageRpc *storagerpc.Client
}

func (its *saveTask) doWorker() {
	batchMsg := make([]*model.Message, 0, batchSize)
	batchM := make([]*nats.Msg, 0, batchSize)
	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()

	for {
		select {
		case m, ok := <-its.messages:
			if !ok {
				if len(batchMsg) > 0 {
					its.saveBatch(batchMsg, batchM)
				}
				return
			}

			msg := &model.Message{}
			err := proto.Unmarshal(m.Data, msg)
			if err != nil {
				_ = m.Ack()
				log.Error("unmarshal message", zap.Error(err))
				continue
			}

			batchMsg = append(batchMsg, msg)
			batchM = append(batchM, m)

			if len(batchMsg) >= batchSize {
				its.saveBatch(batchMsg, batchM)
				batchMsg = batchMsg[:0]
				batchM = batchM[:0]
			}
		case <-timer.C:
			if len(batchMsg) > 0 {
				its.saveBatch(batchMsg, batchM)
				batchMsg = batchMsg[:0]
				batchM = batchM[:0]
			}
		}
	}
}
func (its *saveTask) saveBatch(batchMsg []*model.Message, batchM []*nats.Msg) {
	for i := 0; i < maxRetries; i++ {
		err := its.storageRpc.SaveMessages(batchMsg)
		if err == nil {
			for i := range batchM {
				_ = batchM[i].Ack()
			}
			log.Info("Save messages success", zap.Int("count", len(batchMsg)))
			return
		}
		log.Error("Error saving messages", zap.Error(err))
		time.Sleep(time.Second * time.Duration(i+1))
	}
	for i := range batchM {
		_ = batchM[i].Nak()
	}
}
