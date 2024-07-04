package dispatch

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jedib0t/go-pretty/v6/table"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/mq"
	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/log"
)

var (
	msgTotal        int64
	offlineMsgTotal int64
	onlineMsgTotal  int64
	savedMsgTotal   int64
)

func init() {
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		for {
			select {
			case <-ticker.C:
				{
					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)
					t.AppendHeader(table.Row{"Processed", "Online", "Offline", "Saved", "Goroutines"})
					t.AppendRows([]table.Row{{
						msgTotal,
						onlineMsgTotal,
						offlineMsgTotal,
						savedMsgTotal,
						runtime.NumGoroutine()},
					})
					t.AppendSeparator()
					t.Render()
				}
			}
		}
	}()
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
	if len(devices) == 0 {
		log.Warn("No devices", zap.String("userId", msg.UserId))
		return nil
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
				fmtAddr := strings.Replace(device.GatewayAddress, ".", "-", -1)
				fmtAddr = strings.Replace(fmtAddr, ":", "-", -1)
				err = producer.Publish(fmt.Sprintf(mq.SendMessageSubjectFormat, fmtAddr), body)
				if err != nil {
					log.Error("Error sending message", zap.Error(err))
					continue
				}
				atomic.AddInt64(&onlineMsgTotal, 1)
				log.Debug("Online message", zap.String("gateway", device.GatewayAddress), zap.String("userId", msg.UserId), zap.String("toId", msg.ToId), zap.String("deviceId", device.DeviceId), zap.Int64("seq", msg.SeqId))
			}
		case model.OfflineState:
			{
				offlineMsgCount, err := storageRpc.GetOfflineMessagesCount(msg.UserId, device.DeviceId)
				if err != nil {
					log.Error("Error getting offline messages count", zap.Error(err))
					continue
				}
				atomic.AddInt64(&offlineMsgTotal, offlineMsgCount)
				//TODO Push notification
				log.Debug("Push notification", zap.String("userId", msg.UserId), zap.String("deviceId", device.DeviceId), zap.Int64("count", offlineMsgCount))
			}
		}
	}

	return nil
}
