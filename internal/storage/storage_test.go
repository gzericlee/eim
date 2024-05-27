package storage

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"eim/internal/model"
	"eim/internal/storage/rpc"
)

var rpcClient *rpc.Client

func init() {
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

func TestSaveUser(t *testing.T) {
	for i := 1; i <= 100000; i++ {
		t.Log(rpcClient.SaveUser(&model.User{
			UserId:     fmt.Sprintf("user-%d", i),
			LoginId:    fmt.Sprintf("user-%d", i),
			UserName:   fmt.Sprintf("用户-%d", i),
			Password:   "pass@word1",
			TenantId:   "bingo",
			TenantName: "品高软件",
		}))
	}
}

func TestGetUser(t *testing.T) {
	user, err := rpcClient.GetUser(fmt.Sprintf("user-2"), "bingo")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(user)
	}
}
