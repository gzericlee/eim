package service

import (
	"context"
	"fmt"

	"github.com/gzericlee/eim/internal/database"
	"github.com/gzericlee/eim/internal/redis"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type MessageService struct {
	database     database.IDatabase
	redisManager *redis.Manager
}

func NewMessageService(database database.IDatabase, redisManager *redis.Manager) *MessageService {
	return &MessageService{
		database:     database,
		redisManager: redisManager,
	}
}

func (its *MessageService) InsertMessage(ctx context.Context, args *rpcmodel.MessageArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.InsertMessage(args.Message)
	if err != nil {
		return fmt.Errorf("save message -> %w", err)
	}

	return nil
}

func (its *MessageService) InsertMessages(ctx context.Context, args *rpcmodel.MessagesArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.InsertMessages(args.Messages)
	if err != nil {
		return fmt.Errorf("save messages -> %w", err)
	}

	return nil
}

func (its *MessageService) SaveOfflineMessages(ctx context.Context, args *rpcmodel.MessagesArgs, reply *rpcmodel.EmptyReply) error {
	err := its.redisManager.SaveOfflineMessages(args.Messages, args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("save offline messages -> %w", err)
	}

	return nil
}

func (its *MessageService) RemoveOfflineMessages(ctx context.Context, args *rpcmodel.MessageIdsArgs, reply *rpcmodel.EmptyReply) error {
	err := its.redisManager.RemoveOfflineMessages(args.MessageIds, args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("remove offline messages -> %w", err)
	}

	return nil
}

func (its *MessageService) GetOfflineMessagesCount(ctx context.Context, args *rpcmodel.OfflineMessagesArgs, reply *rpcmodel.OfflineMessagesCountReply) error {
	count, err := its.redisManager.GetOfflineMessagesCount(args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("get offline messages count -> %w", err)
	}

	reply.Count = count

	return nil
}

func (its *MessageService) GetOfflineMessages(ctx context.Context, args *rpcmodel.OfflineMessagesArgs, reply *rpcmodel.MessagesReply) error {
	result, err := its.redisManager.GetOfflineMessages(args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("get offline messages -> %w", err)
	}

	reply.Messages = result

	return nil
}
