package dispatch

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/util/log"

	"eim/internal/mq"
	"eim/internal/redis"
)

func toGroup(msg *model.Message, redisManager *redis.Manager, producer mq.Producer) error {
	members, err := redisManager.GetBizMembers(model.BizGroup, msg.ToId)
	if err != nil {
		return fmt.Errorf("get group members -> %w", err)
	}
	members = append(members, msg.FromId)
	for _, userId := range members {
		msg.UserId = userId
		err = toUser(msg, redisManager, producer)
		if err != nil {
			return fmt.Errorf("send message to user -> %w", err)
		}
	}

	return nil
}

func toUser(msg *model.Message, redisManager *redis.Manager, producer mq.Producer) error {
	devices, err := redisManager.GetUserDevices(msg.UserId)
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
		key := fmt.Sprintf("%v:offline:%v:%v", msg.UserId, msg.ToId, device.DeviceId)
		offlineCount, err := redisManager.Incr(key)
		if err != nil {
			return fmt.Errorf("incr offline count -> %w", err)
		}
		if device.State == model.OnlineState {
			err = producer.Publish(fmt.Sprintf(mq.MessageSendSubject, strings.Replace(device.GatewayIp, ".", "-", -1)), body)
			if err != nil {
				return fmt.Errorf("publish message -> %w", err)
			}
			log.Debug("Online message", zap.String("gateway", device.GatewayIp), zap.String("userId", msg.UserId), zap.String("toId", msg.ToId), zap.String("deviceId", device.DeviceId), zap.Int64("seq", msg.SeqId))
		} else {
			log.Debug("Offline message", zap.Int64("count", offlineCount), zap.String("userId", msg.UserId), zap.String("toId", msg.ToId), zap.String("deviceId", device.DeviceId))
		}
	}

	return nil
}
