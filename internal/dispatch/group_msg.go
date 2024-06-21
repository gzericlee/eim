package dispatch

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"

	"eim/internal/model"
	"eim/internal/mq"
	storagerpc "eim/internal/storage/rpc"
	"eim/util/log"
)

type GroupMessageHandler struct {
	storageRpc *storagerpc.Client
	producer   mq.Producer
}

func NewGroupMessageHandler(storageRpc *storagerpc.Client, producer mq.Producer) *GroupMessageHandler {
	return &GroupMessageHandler{
		storageRpc: storageRpc,
		producer:   producer,
	}
}

func (its *GroupMessageHandler) HandleMessage(m *nats.Msg) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	if m.Data == nil || len(m.Data) == 0 {
		return m.Ack()
	}

	msg := &model.Message{}
	err := proto.Unmarshal(m.Data, msg)
	if err != nil {
		return fmt.Errorf("unmarshal message -> %w", err)
	}

	err = toGroup(msg, its.storageRpc, its.producer)
	if err != nil {
		return fmt.Errorf("send message to group -> %w", err)
	}

	return m.Ack()
}
