package storage

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"eim/global"
	"eim/internal/database/maindb"
	"eim/internal/redis"
	"eim/model"
	"eim/proto/pb"
)

var rpcClient *RpcClient

func init() {
	global.InitLogger()
	fmt.Println(maindb.InitDBEngine("mysql", "root:pass@word1@tcp(10.8.12.23:4000)/eim?charset=utf8mb4&parseTime=True&loc=Local"))
	go func() {
		fmt.Println(InitStorageServer("0.0.0.0", 10000, []string{"10.8.12.23:2379", "10.8.12.23:2479", "10.8.12.23:2579"}))
	}()

	err := redis.InitRedisClusterClient([]string{"10.8.12.23:7001", "10.8.12.23:7002", "10.8.12.23:7003"}, "pass@word1")
	if err != nil {
		panic(err)
	}

	rpcClient, err = NewRpcClient([]string{"10.8.12.23:2379", "10.8.12.23:2479", "10.8.12.23:2579"})
	if err != nil {
		panic(err)
	}
}

func TestMessage_Save(t *testing.T) {
	s := 0
	for i := 0; i < 150000; i++ {
		msg := &pb.Message{
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
		device := &model.Device{
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
