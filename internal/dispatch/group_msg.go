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
	task       *saveTask
	storageRpc *storagerpc.Client
	producer   mq.Producer
}

func NewGroupMessageHandler(storageRpc *storagerpc.Client, producer mq.Producer) *GroupMessageHandler {
	task := &saveTask{messages: make(chan *nats.Msg, 1000), storageRpc: storageRpc}
	go task.doWorker()
	return &GroupMessageHandler{
		task:       task,
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
		_ = m.Ack()
		return fmt.Errorf("message data is nil")
	}

	its.task.messages <- m

	msg := &model.Message{}
	err := proto.Unmarshal(m.Data, msg)
	if err != nil {
		return fmt.Errorf("unmarshal message -> %w", err)
	}

	err = toGroup(msg, its.storageRpc, its.producer)
	if err != nil {
		return fmt.Errorf("send message to group -> %w", err)
	}

	return nil
}
