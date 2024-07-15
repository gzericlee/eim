package model

import "github.com/gzericlee/eim/internal/model"

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
