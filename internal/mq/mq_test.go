package mq

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"

	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/pkg/snowflake"
)

var producer IProducer
var consumer IConsumer

func init() {
	snowflake.NewGenerator(snowflake.GeneratorConfig{
		RedisEndpoints: []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
		RedisPassword:  "pass@word1",
		MaxWorkerId:    1000,
		MinWorkerId:    0,
		NodeCount:      1,
	})

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
			MsgId:      1,
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

func (h *testHandler) Process(m *nats.Msg) error {
	msg := &model.Message{}
	err := proto.Unmarshal(m.Data, msg)
	return err
}

func TestSubscribe(t *testing.T) {
	err := consumer.Subscribe(fmt.Sprintf(SendMessageSubject, "192_168_3_58"), "192_168_3_58", &testHandler{})
	if err != nil {
		return
	}
	time.Sleep(time.Second * 5)
}
