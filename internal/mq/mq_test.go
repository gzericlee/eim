package mq

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"eim/internal/model"
)

var producer Producer

func init() {
	var err error
	producer, err = NewProducer([]string{"127.0.0.1:4151", "127.0.0.1:4161", "127.0.0.1:4171"})
	if err != nil {
		panic(err)
	}
}

func BenchmarkPool_Publish(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
		body, _ := msg.Serialize()
		err := producer.Publish("test", body)
		if err != nil {
			b.Log(err)
			return
		}
	}
}
