package dispatch

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"eim/internal/model"
	"eim/internal/mq"
	storagerpc "eim/internal/storage/rpc"
	"eim/util/log"
)

type UserMessageHandler struct {
	StorageRpc *storagerpc.Client
	Producer   mq.Producer
}

func (its *UserMessageHandler) HandleMessage(data []byte) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	if data == nil || len(data) == 0 {
		return nil
	}

	msg := &model.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("unmarshal message -> %w", err)
	}

	err = its.StorageRpc.SaveMessage(msg)
	if err != nil {
		return fmt.Errorf("save message -> %w", err)
	}

	//发送者多端同步
	msg.UserId = msg.FromId
	err = toUser(msg, its.StorageRpc, its.Producer)
	if err != nil {
		return fmt.Errorf("dispatch user message to send user -> %w", err)
	}

	//接收者多端推送
	msg.UserId = msg.ToId
	err = toUser(msg, its.StorageRpc, its.Producer)
	if err != nil {
		return fmt.Errorf("dispatch user message to receive user -> %w", err)
	}

	return nil
}
