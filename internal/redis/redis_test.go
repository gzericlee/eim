package redis

import (
	"context"
	"testing"
	"time"

	"eim/internal/model"
)

var manager *Manager

func init() {
	var err error
	manager, err = NewManager(Config{
		RedisEndpoints:       []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
		RedisPassword:        "pass@word1",
		OfflineMessageExpire: time.Minute * 2,
		OfflineDeviceExpire:  time.Minute * 2,
	})
	if err != nil {
		panic(err)
	}
}

func TestManager_GetGateways(t *testing.T) {
	gateways, err := manager.GetGateways()
	if err != nil {
		t.Error(err)
		return
	}
	for _, gateway := range gateways {
		t.Log(gateway.Ip)
	}
}

func TestManager_SaveOfflineMessages(t *testing.T) {
	var msgs []*model.Message
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &model.Message{
			MsgId:    int64(i),
			FromId:   "user-1",
			ToId:     "user-2",
			SeqId:    1,
			Content:  "Hello, World!",
			SendTime: time.Now().Unix(),
		})
	}
	t.Log(manager.SaveOfflineMessages(msgs, "user-1", "device-1"))
	t.Log(manager.GetOfflineMessagesCount("user-1", "device-1"))
	t.Log(manager.GetOfflineMessages("user-1", "device-1"))
}

func TestManager_RemoveAllByKeys(t *testing.T) {
	keys, _ := manager.getAllKeys("members:*")
	for _, key := range keys {
		err := manager.redisClient.Del(context.Background(), key).Err()
		if err != nil {
			t.Error(err)
			return
		}
	}
}
