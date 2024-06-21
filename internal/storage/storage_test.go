package storage

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"

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
		err := rpcClient.SaveMessages([]*model.Message{msg})
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
	for i := 1; i <= 10; i++ {
		t.Log(rpcClient.SaveBiz(&model.Biz{
			BizId:      fmt.Sprintf("user-%d", i),
			BizType:    model.Biz_USER,
			BizName:    fmt.Sprintf("用户-%d", i),
			TenantId:   "bingo",
			TenantName: "品高软件",
			Attributes: map[string]*anypb.Any{},
		}))
	}
	for i := 1; i <= 10; i++ {
		user, err := rpcClient.GetBiz(fmt.Sprintf("user-%d", i), "bingo")
		if err != nil {
			t.Log(err)
		} else {
			t.Log(user)
		}
	}
}

func TestGetUser(t *testing.T) {
	for i := 1; i <= 10; i++ {
		user, err := rpcClient.GetBiz(fmt.Sprintf("user-%d", i), "bingo")
		if err != nil {
			t.Log(err)
		} else {
			t.Log(user)
		}
	}
}
