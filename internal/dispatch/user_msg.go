package dispatch

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/mq"
	storagerpc "eim/internal/storage/rpc"
	"eim/util/log"
)

type UserMessageHandler struct {
	task       *saveTask
	storageRpc *storagerpc.Client
	producer   mq.Producer
}

func NewUserMessageHandler(storageRpc *storagerpc.Client, producer mq.Producer) *UserMessageHandler {
	task := &saveTask{messages: make(chan *nats.Msg, 1000), storageRpc: storageRpc}
	go task.doWorker()
	return &UserMessageHandler{
		task:       task,
		storageRpc: storageRpc,
		producer:   producer,
	}
}

func (its *UserMessageHandler) HandleMessage(m *nats.Msg) error {
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
		_ = m.Nak()
		log.Error("unmarshal message", zap.Error(err))
	}

	//发送者多端同步
	msg.UserId = msg.FromId
	err = toUser(msg, its.storageRpc, its.producer)
	if err != nil {
		return fmt.Errorf("dispatch user message to send user -> %w", err)
	}

	//接收者多端推送
	msg.UserId = msg.ToId
	err = toUser(msg, its.storageRpc, its.producer)
	if err != nil {
		return fmt.Errorf("dispatch user message to receive user -> %w", err)
	}

	return nil
}
