package dispatch

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"

	"eim/internal/model"
	storagerpc "eim/internal/storage/rpc"
)

type SaveMessageHandler struct {
	storageRpc *storagerpc.Client
}

func NewSaveMessageHandler(storageRpc *storagerpc.Client) *SaveMessageHandler {
	return &SaveMessageHandler{
		storageRpc: storageRpc,
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

	err = its.storageRpc.InsertMessage(msg)
	if err != nil {
		return fmt.Errorf("insert message -> %w", err)
	}

	savedMsgTotal.Add(1)

	return m.Ack()
}
