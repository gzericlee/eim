package dispatch

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"

	"eim/internal/mq"
	"eim/internal/redis"
	"eim/pkg/log"
)

func toGroup(msg *model.Message, redisManager *redis.Manager, producer mq.Producer) error {
	members, err := redisManager.GetBizMembers(model.BizGroup, msg.ToId)
	if err != nil {
		log.Error("Error getting group members", zap.String("groupId", msg.ToId), zap.Error(err))
		return err
	}
	members = append(members, msg.FromId)
	for _, userId := range members {
		msg.UserId = userId
		err = toUser(msg, redisManager, producer)
		if err != nil {
			return err
		}
	}

	return nil
}

func toUser(msg *model.Message, redisManager *redis.Manager, producer mq.Producer) error {
	devices, err := redisManager.GetUserDevices(msg.UserId)
	if err != nil {
		log.Error("Error getting from user devices", zap.String("userId", msg.FromId), zap.Error(err))
		return err
	}
	body, err := proto.Marshal(msg)
	if err != nil {
		log.Error("Error marshalling message", zap.Error(err))
		return err
	}

	for _, device := range devices {
		//if msg.FromDevice == device.DeviceId {
		//	continue
		//}
		key := fmt.Sprintf("%v:offline:%v:%v", msg.UserId, msg.ToId, device.DeviceId)
		offlineCount, err := redisManager.Incr(key)
		if err != nil {
			log.Error("Error increasing offline count", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId), zap.Error(err))
			return err
		}
		if device.State == model.OnlineState {
			err = producer.Publish(fmt.Sprintf(mq.MessageSendSubject, strings.Replace(device.GatewayIp, ".", "-", -1)), body)
			if err != nil {
				log.Error("Error publishing message", zap.Error(err))
				return err
			}
			log.Debug("Online message", zap.String("gateway", device.GatewayIp), zap.String("userId", msg.UserId), zap.String("toId", msg.ToId), zap.String("deviceId", device.DeviceId), zap.Int64("seq", msg.SeqId))
		} else {
			log.Debug("Offline message", zap.Int64("count", offlineCount), zap.String("userId", msg.UserId), zap.String("toId", msg.ToId), zap.String("deviceId", device.DeviceId))
		}
	}

	return nil
}
