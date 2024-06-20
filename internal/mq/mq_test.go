package mq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"

	"eim/internal/model"
	"eim/pkg/snowflake"
)

var producer Producer
var consumer Consumer

func init() {
	snowflake.Init([]string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"}, "pass@word1")

	var err error
	producer, err = NewProducer([]string{"127.0.0.1:4222"})
	if err != nil {
		panic(err)
	}

	consumer, err = NewConsumer([]string{"127.0.0.1:4222"})
	if err != nil {
		panic(err)
	}
}

func BenchmarkPublish(b *testing.B) {
	for i := 0; i < 1; i++ {
		msg := &model.Message{
			MsgId:      uuid.New().String(),
			SeqId:      1,
			MsgType:    1,
			Content:    time.Now().String(),
			FromType:   1,
			FromId:     "1",
			FromDevice: "1",
			ToType:     1,
			ToId:       "1",
			ToDevice:   "1",
			SendTime:   time.Now().UnixMilli(),
		}
		body, _ := proto.Marshal(msg)
		err := producer.Publish(fmt.Sprintf(SendMessageSubject, "192_168_3_58"), body)
		if err != nil {
			b.Log(err)
			return
		}
	}
}

type testHandler struct{}

func (h *testHandler) HandleMessage(data []byte) error {
	msg := &model.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	log.Printf("%s", data)
	return nil
}

func TestSubscribe(t *testing.T) {
	err := consumer.Subscribe(fmt.Sprintf(SendMessageSubject, "192_168_3_58"), "192_168_3_58", &testHandler{})
	if err != nil {
		return
	}
	time.Sleep(time.Second * 5)
}
