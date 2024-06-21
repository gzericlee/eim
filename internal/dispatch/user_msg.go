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

type UserMessageHandler struct {
	storageRpc *storagerpc.Client
	producer   mq.Producer
}

func NewUserMessageHandler(storageRpc *storagerpc.Client, producer mq.Producer) *UserMessageHandler {
	return &UserMessageHandler{
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
		return m.Ack()
	}

	msg := &model.Message{}
	err := proto.Unmarshal(m.Data, msg)
	if err != nil {
		return m.Ack()
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

	return m.Ack()
}
