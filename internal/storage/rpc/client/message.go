package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type MessageClient struct {
	*rpcxclient.XClientPool
}

func (its *MessageClient) InsertMessage(message *model.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "InsertMessage", &rpcmodel.MessageArgs{Message: message}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertMessages -> %w", err)
	}
	return nil
}

func (its *MessageClient) InsertMessages(messages []*model.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "InsertMessages", &rpcmodel.MessagesArgs{Messages: messages}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertMessages -> %w", err)
	}
	return nil
}

func (its *MessageClient) SaveOfflineMessages(msgs []*model.Message, userId, deviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "SaveOfflineMessages", &rpcmodel.MessagesArgs{Messages: msgs, UserId: userId, DeviceId: deviceId}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call SaveOfflineMessages -> %w", err)
	}
	return nil
}

func (its *MessageClient) RemoveOfflineMessages(msgIds []string, userId, deviceId string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.Get().Call(ctx, "RemoveOfflineMessages", &rpcmodel.MessageIdsArgs{MessageIds: msgIds, UserId: userId, DeviceId: deviceId}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RemoveOfflineMessages -> %w", err)
	}
	return nil
}

func (its *MessageClient) GetOfflineMessagesCount(userId, deviceId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.OfflineMessagesCountReply{}
	err := its.Get().Call(ctx, "GetOfflineMessagesCount", &rpcmodel.OfflineMessagesArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return 0, fmt.Errorf("call GetOfflineMessagesCount -> %w", err)
	}

	return reply.Count, nil
}

func (its *MessageClient) GetOfflineMessages(userId, deviceId string) ([]*model.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.MessagesReply{}
	err := its.Get().Call(ctx, "GetOfflineMessages", &rpcmodel.OfflineMessagesArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetOfflineMessages -> %w", err)
	}

	return reply.Messages, nil
}
