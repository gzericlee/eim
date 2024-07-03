package dispatch

import (
	"fmt"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"

	"eim/internal/model"
	"eim/internal/mq"
	storagerpc "eim/internal/storage/rpc"
)

type GroupMessageHandler struct {
	storageRpc *storagerpc.Client
	producer   mq.IProducer
}

func NewGroupMessageHandler(storageRpc *storagerpc.Client, producer mq.IProducer) *GroupMessageHandler {
	return &GroupMessageHandler{
		storageRpc: storageRpc,
		producer:   producer,
	}
}

func (its *GroupMessageHandler) Process(m *nats.Msg) error {
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

	atomic.AddInt64(&groupMsgTotal, 1)

	return m.Ack()
}
