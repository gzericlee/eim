package dispatch

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/internal/model/consts"
	"github.com/gzericlee/eim/internal/mq"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
	"github.com/gzericlee/eim/pkg/log"
)

type UserMessageHandler struct {
	bizMemberRpc *storagerpc.BizMemberClient
	messageRpc   *storagerpc.MessageClient
	deviceRpc    *storagerpc.DeviceClient
	producer     mq.IProducer
}

func NewUserMessageHandler(bizMemberRpc *storagerpc.BizMemberClient, messageRpc *storagerpc.MessageClient, deviceRpc *storagerpc.DeviceClient, producer mq.IProducer) *UserMessageHandler {
	return &UserMessageHandler{
		bizMemberRpc: bizMemberRpc,
		messageRpc:   messageRpc,
		deviceRpc:    deviceRpc,
		producer:     producer,
	}
}

func (its *UserMessageHandler) Process(m *nats.Msg) error {
	if m.Data == nil || len(m.Data) == 0 {
		return m.Ack()
	}

	msg := &model.Message{}
	err := proto.Unmarshal(m.Data, msg)
	if err != nil {
		return m.Ack()
	}

	msg.UserId = msg.FromId
	msg.TenantId = msg.FromTenant
	err = its.publish(*msg)
	if err != nil {
		return fmt.Errorf("dispatch user message to send user -> %w", err)
	}

	msg.UserId = msg.ToId
	msg.TenantId = msg.ToTenant
	err = its.publish(*msg)
	if err != nil {
		return fmt.Errorf("dispatch user message to receive user -> %w", err)
	}

	msgTotal.Add(1)

	return m.Ack()
}

func (its *UserMessageHandler) publish(msg model.Message) error {
	devices, err := its.deviceRpc.GetDevices(msg.UserId, msg.TenantId)
	if err != nil {
		return fmt.Errorf("get user devices -> %w", err)
	}
	if len(devices) == 0 {
		log.Warn("No devices", zap.String("userId", msg.UserId))
		return nil
	}

	body, err := proto.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("marshal message -> %w", err)
	}

	for _, device := range devices {
		if msg.FromDevice == device.DeviceId {
			continue
		}

		err = its.messageRpc.SaveOfflineMessages([]*model.Message{&msg}, msg.UserId, device.DeviceId)
		if err != nil {
			log.Error("Error saving offline messages", zap.Error(err))
			continue
		}

		switch device.State {
		case consts.StatusOnline:
			{
				fmtAddr := strings.Replace(device.GatewayAddr, ".", "-", -1)
				fmtAddr = strings.Replace(fmtAddr, ":", "-", -1)
				err = its.producer.Publish(fmt.Sprintf(mq.SendMessageSubjectFormat, fmtAddr), body)
				if err != nil {
					log.Error("Error sending message", zap.Error(err))
					continue
				}
				onlineMsgTotal.Add(1)
				log.Debug("Online message", zap.String("gateway", device.GatewayAddr), zap.String("userId", msg.UserId), zap.String("toId", msg.ToId), zap.String("deviceId", device.DeviceId), zap.Int64("seq", msg.SeqId))
			}
		case consts.StatusOffline:
			{
				offlineMsgCount, err := its.messageRpc.GetOfflineMessagesCount(msg.UserId, device.DeviceId)
				if err != nil {
					log.Error("Error getting offline messages count", zap.Error(err))
					continue
				}
				offlineMsgTotal.Add(offlineMsgCount)
				//TODO Push notification
				log.Info("Push notification", zap.String("userId", msg.UserId), zap.String("deviceId", device.DeviceId), zap.Int64("count", offlineMsgCount))
			}
		}
	}

	return nil
}
