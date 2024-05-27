package rpc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/redis"
	"eim/util/log"

	"eim/internal/database"
)

type Message struct {
	database     database.IDatabase
	redisManager *redis.Manager
}

type MessageArgs struct {
	Message *model.Message
}

type MessageIdsArgs struct {
	MsgIds   []interface{}
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

type OfflineMessagesReply struct {
	Messages []*model.Message
}

func (its *Message) SaveMessage(ctx context.Context, args *MessageArgs, reply *EmptyReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.database.SaveMessage(args.Message)
	if err != nil {
		return fmt.Errorf("save message -> %w", err)
	}

	return nil
}

func (its *Message) SaveOfflineMessageIds(ctx context.Context, args *MessageIdsArgs, reply *EmptyReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.redisManager.SaveOfflineMessageIds(args.MsgIds, args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("save offline message ids -> %w", err)
	}

	return nil
}

func (its *Message) RemoveOfflineMessageIds(ctx context.Context, args *MessageIdsArgs, reply *EmptyReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.redisManager.RemoveOfflineMessageIds(args.MsgIds, args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("remove offline message ids -> %w", err)
	}

	return nil
}

func (its *Message) GetOfflineMessagesCount(ctx context.Context, args *OfflineMessagesArgs, reply *OfflineMessagesCountReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	count, err := its.redisManager.GetOfflineMessagesCount(args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("get offline messages count -> %w", err)
	}

	reply.Count = count

	return nil
}

func (its *Message) GetOfflineMessages(ctx context.Context, args *OfflineMessagesArgs, reply *OfflineMessagesReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	result, err := its.redisManager.GetOfflineMessages(args.UserId, args.DeviceId)
	if err != nil {
		return fmt.Errorf("get offline messages -> %w", err)
	}

	var allMsgIds []int64
	for _, msgIds := range result {
		for _, msgId := range msgIds {
			intId, err := strconv.ParseInt(msgId, 10, 64)
			if err != nil {
				log.Error("Error converting message id to int", zap.String("msgId", msgId), zap.Error(err))
				continue
			}
			allMsgIds = append(allMsgIds, intId)
		}
	}

	if len(allMsgIds) == 0 {
		return nil
	}

	messages, err := its.database.GetMessagesByIds(allMsgIds)
	if err != nil {
		return fmt.Errorf("get messages by ids -> %w", err)
	}

	reply.Messages = messages

	return nil
}
