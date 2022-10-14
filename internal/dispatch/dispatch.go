package dispatch

import (
	"fmt"

	"go.uber.org/zap"

	"eim/internal/nsq/producer"
	"eim/internal/redis"
	"eim/internal/types"
	"eim/pkg/log"
)

func toUser(msg *types.Message) error {
	//获取用户设备
	devices, err := redis.GetUserDevices(msg.UserId)
	if err != nil {
		log.Error("Error getting User devices", zap.String("userId", msg.UserId), zap.Error(err))
		return err
	}
	body, _ := msg.Serialize()

	for _, device := range devices {
		if msg.FromDevice == device.DeviceId {
			continue
		}
		key := fmt.Sprintf("%v:offline:%v:%v", msg.UserId, msg.ToId, device.DeviceId)
		offlineCount, err := redis.Incr(key)
		if err != nil {
			log.Error("Error increasing offline count", zap.String("userId", msg.UserId), zap.String("deviceId", device.DeviceId), zap.Error(err))
			return err
		}
		if device.State == types.OnlineState {
			err = producer.PublishAsync(device.GatewayIp+"_send", body)
			if err != nil {
				log.Error("Error publishing message", zap.Error(err))
				return err
			}
			log.Debug("Online message", zap.String("gateway", device.GatewayIp), zap.String("userId", msg.UserId), zap.String("deviceId", device.DeviceId))
		} else {
			log.Debug("Offline message", zap.Int64("count", offlineCount), zap.String("userId", msg.UserId), zap.String("deviceId", device.DeviceId))
		}
	}

	return nil
}
