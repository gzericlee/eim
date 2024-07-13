package dispatch

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"

	"eim/internal/model"
	"eim/internal/mq"
	storagerpc "eim/internal/storage/rpc"
)

type GroupMessageHandler struct {
	storageRpc         *storagerpc.Client
	producer           mq.IProducer
	userMessageHandler *UserMessageHandler
}

func NewGroupMessageHandler(storageRpc *storagerpc.Client, producer mq.IProducer) *GroupMessageHandler {
	return &GroupMessageHandler{
		storageRpc:         storageRpc,
		producer:           producer,
		userMessageHandler: NewUserMessageHandler(storageRpc, producer),
	}
}

func (its *GroupMessageHandler) Process(m *nats.Msg) error {
	if m.Data == nil || len(m.Data) == 0 {
		return m.Ack()
	}

	msg := &model.Message{}
	err := proto.Unmarshal(m.Data, msg)
	if err != nil {
		return fmt.Errorf("unmarshal message -> %w", err)
	}

	err = its.publish(msg)
	if err != nil {
		return fmt.Errorf("send message to group -> %w", err)
	}

	msgTotal.Add(1)

	return m.Ack()
}

func (its *GroupMessageHandler) publish(msg *model.Message) error {
	// 获取群组成员,ToId，ToTenantId，是指群组ID和群组的租户ID
	members, err := its.storageRpc.GetBizMembers(msg.ToId, msg.ToTenant)
	if err != nil {
		return fmt.Errorf("get group members -> %w", err)
	}

	for _, member := range members {
		// member格式为userId@tenantId，群组成员可能是其他租户的用户
		userId := strings.Split(member, "@")[0]
		tenantId := strings.Split(member, "@")[1]
		msg.UserId = userId
		msg.TenantId = tenantId
		err = its.userMessageHandler.publish(*msg)
		if err != nil {
			return fmt.Errorf("send message to user -> %w", err)
		}
	}

	return nil
}
