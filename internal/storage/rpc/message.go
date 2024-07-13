package rpc

import (
	"context"
	"fmt"

	"eim/internal/database"
	"eim/internal/model"
	"eim/internal/redis"
)

type Message struct {
	database     database.IDatabase
	redisManager *redis.Manager
}

type MessageIdsArgs struct {
	MessageIds []string
	UserId     string
	DeviceId   string
}

type MessageArgs struct {
	Message *model.Message
}

type MessagesArgs struct {
	Messages []*model.Message
	UserId   string
	DeviceId string
}

type OfflineMessagesArgs struct {
	UserId   string
	DeviceId string
}

type OfflineMessagesCountReply struct {
	Count int64
}

type MessagesReply struct {
	Messages []*model.Message
}

func (its *Message) InsertMessage(ctx context.Context, args *MessageArgs, reply *EmptyReply) error {
	err := its.database.InsertMessage(args.Message)
	if err != nil {
		return fmt.Errorf("save message -> %w", err)
	}

	return nil
}

func (its *Message) InsertMessages(ctx context.Context, args *MessagesArgs, reply *EmptyReply) error {
	err := its.database.InsertMessages(args.Messages)
	if err != nil {
		return fmt.Errorf("save messages -> %w", err)
	}

	return nil
}

func (its *Message) SaveOfflineMessages(ctx context.Context, args *MessagesArgs, reply *EmptyReply) error {
	err := its.redisManager.SaveOfflineMessages(args.Messages, args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("save offline messages -> %w", err)
	}

	return nil
}

func (its *Message) RemoveOfflineMessages(ctx context.Context, args *MessageIdsArgs, reply *EmptyReply) error {
	err := its.redisManager.RemoveOfflineMessages(args.MessageIds, args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("remove offline messages -> %w", err)
	}

	return nil
}

func (its *Message) GetOfflineMessagesCount(ctx context.Context, args *OfflineMessagesArgs, reply *OfflineMessagesCountReply) error {
	count, err := its.redisManager.GetOfflineMessagesCount(args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("get offline messages count -> %w", err)
	}

	reply.Count = count

	return nil
}

func (its *Message) GetOfflineMessages(ctx context.Context, args *OfflineMessagesArgs, reply *MessagesReply) error {
	result, err := its.redisManager.GetOfflineMessages(args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("get offline messages -> %w", err)
	}

	reply.Messages = result

	return nil
}
