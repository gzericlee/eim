package storage

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"eim/internal/database/maindb"
	"eim/internal/redis"
	"eim/internal/storage/rpc"
	"eim/internal/types"
)

var rpcClient *rpc.Client

func init() {
	fmt.Println(maindb.InitDBEngine("mysql", "root:pass@word1@tcp(127.0.0.1:4000)/eim?charset=utf8mb4&parseTime=True&loc=Local"))
	go func() {
		fmt.Println(InitStorageServer("0.0.0.0", 10000, []string{"127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"}))
	}()

	err := redis.InitRedisClusterClient([]string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003"}, "pass@word1")
	if err != nil {
		panic(err)
	}

	rpcClient, err = rpc.NewClient([]string{"127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"})
	if err != nil {
		panic(err)
	}
}

func TestMessage_Save(t *testing.T) {
	s := 0
	for i := 0; i < 150000; i++ {
		msg := &types.Message{
			MsgId:      uuid.New().String(),
			SeqId:      1,
			MsgType:    1,
			Content:    "",
			FromType:   1,
			FromId:     uuid.New().String(),
			FromDevice: uuid.New().String(),
			ToType:     0,
			ToId:       uuid.New().String(),
			ToDevice:   uuid.New().String(),
			SendTime:   0,
		}
		err := rpcClient.SaveMessage(msg)
		if err != nil {
			t.Log(err)
		} else {
			s++
		}
	}
	t.Log(s)
}

func TestDevice_Save(t *testing.T) {
	s := 0
	for i := 0; i < 150000; i++ {
		device := &types.Device{
			DeviceId:      uuid.New().String(),
			UserId:        uuid.New().String(),
			DeviceType:    uuid.New().String(),
			DeviceVersion: uuid.New().String(),
			GatewayIp:     uuid.New().String(),
			OnlineAt:      nil,
			OfflineAt:     nil,
			State:         0,
		}
		err := rpcClient.SaveDevice(device)
		if err != nil {
			t.Log(err)
		} else {
			s++
		}
	}
	t.Log(s)
}
