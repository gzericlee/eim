package producer

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"

	"eim/model"
)

func init() {
	log.Println(InitProducers([]string{"10.8.12.23:4161", "10.8.12.23:4261"}))
}

func BenchmarkPool_Publish(b *testing.B) {
	wg := sync.WaitGroup{}
	doneChan := make(chan *nsq.ProducerTransaction, 1)
	go func() {
		for _ = range doneChan {
			wg.Done()
		}
	}()
	b.N = 100000
	wg.Add(b.N)
	time.Sleep(time.Second)
	for i := 0; i < b.N; i++ {
		msg := &model.Message{
			MsgId:      uuid.New().String(),
			SeqId:      1,
			MsgType:    1,
			Content:    time.Now().String(),
			FromType:   1,
			FromId:     "1",
			FromName:   "1",
			FromDevice: "1",
			ToType:     1,
			ToId:       "1",
			ToName:     "1",
			ToDevice:   "1",
			SendTime:   time.Now().UnixMilli(),
		}
		body, _ := msg.Serialize()
		err := pool.PublishAsync("test", body, doneChan, msg.MsgId)
		if err != nil {
			b.Log(err)
			return
		}
	}
	wg.Wait()
}
