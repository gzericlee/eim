package dispatch

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/ants"
	"go.uber.org/zap"

	"eim/internal/model"
	storagerpc "eim/internal/storage/rpc"
	"eim/util/log"

	"eim/internal/mq"
)

var taskPool *ants.Pool

func init() {
	var err error
	taskPool, err = ants.NewPoolPreMalloc(2048)
	if err != nil {
		panic(fmt.Sprintf("new ants pool -> %v", err))
	}
}

func toGroup(msg *model.Message, storageRpc *storagerpc.Client, producer mq.IProducer) error {
	members, err := storageRpc.GetBizMembers(model.BizGroup, msg.ToId)
	if err != nil {
		return fmt.Errorf("get group members -> %w", err)
	}

	members = append(members, msg.FromId)
	for _, userId := range members {
		msg.UserId = userId
		err = toUser(msg, storageRpc, producer)
		if err != nil {
			return fmt.Errorf("send message to user -> %w", err)
		}
	}

	return nil
}

func toUser(msg *model.Message, storageRpc *storagerpc.Client, producer mq.IProducer) error {
	devices, err := storageRpc.GetDevices(msg.UserId)
	if err != nil {
		return fmt.Errorf("get user devices -> %w", err)
	}

	body, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message -> %w", err)
	}

	for _, device := range devices {
		if msg.FromDevice == device.DeviceId {
			continue
		}

		err = storageRpc.SaveOfflineMessages([]*model.Message{msg}, msg.UserId, device.DeviceId)
		if err != nil {
			log.Error("Error saving offline messages", zap.Error(err))
			continue
		}

		switch device.State {
		case model.OnlineState:
			{
				err = producer.Publish(fmt.Sprintf(mq.SendMessageSubject, strings.Replace(device.GatewayIp, ".", "-", -1)), body)
				if err != nil {
					log.Error("Error sending message", zap.Error(err))
					continue
				}

				log.Debug("Online message", zap.String("gateway", device.GatewayIp), zap.String("userId", msg.UserId), zap.String("toId", msg.ToId), zap.String("deviceId", device.DeviceId), zap.Int64("seq", msg.SeqId))
			}
		case model.OfflineState:
			{
				offlineMsgCount, err := storageRpc.GetOfflineMessagesCount(msg.UserId, device.DeviceId)
				if err != nil {
					log.Error("Error getting offline messages count", zap.Error(err))
					continue
				}

				//TODO Push notification

				log.Debug("Push notification", zap.String("userId", msg.UserId), zap.String("deviceId", device.DeviceId), zap.Int64("count", offlineMsgCount))
			}
		}
	}

	return nil
}
