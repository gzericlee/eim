package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"eim/internal/database"
	"eim/internal/model"
	"eim/internal/storage/rpc"
)

var rpcClient *rpc.Client

func init() {
	go func() {
		fmt.Println(rpc.StartServer(rpc.Config{Ip: "0.0.0.0", Port: 10000, EtcdEndpoints: []string{"127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"}, DatabaseConnection: "mongodb://admin:pass@word1@127.0.0.1:27017", DatabaseDriver: database.DriverMongoDB}))
	}()

	time.Sleep(time.Second)

	var err error
	rpcClient, err = rpc.NewClient([]string{"127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"})
	if err != nil {
		panic(err)
	}
}

func TestMessage_Save(t *testing.T) {
	s := 0
	for i := 0; i < 150000; i++ {
		msg := &model.Message{
			MsgId:      0,
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
