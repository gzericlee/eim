package dispatch

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"

	"github.com/gzericlee/eim/internal/model"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

type SaveMessageHandler struct {
	messageRpc *storagerpc.MessageClient
}

func NewSaveMessageHandler(messageRpc *storagerpc.MessageClient) *SaveMessageHandler {
	return &SaveMessageHandler{
		messageRpc: messageRpc,
	}
}

func (its *SaveMessageHandler) Process(m *nats.Msg) error {
	if m.Data == nil || len(m.Data) == 0 {
		return m.Ack()
	}

	msg := &model.Message{}
	err := proto.Unmarshal(m.Data, msg)
	if err != nil {
		return m.Ack()
	}

	err = its.messageRpc.InsertMessage(msg)
	if err != nil {
		return fmt.Errorf("insert message -> %w", err)
	}

	savedMsgTotal.Add(1)

	return m.Ack()
}
