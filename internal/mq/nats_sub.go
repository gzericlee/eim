package mq

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nats-io/nats.go"
	"github.com/panjf2000/ants"
	"go.uber.org/zap"

	"eim/util/log"
)

type natsConsumer struct {
	conn      *nats.Conn
	jsContext nats.JetStreamContext
	taskPool  *ants.Pool
}

func newNatsConsumer(endpoints []string) (IConsumer, error) {
	conn, err := nats.Connect(
		strings.Join(endpoints, ","),
		nats.Name("eim"),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(5),
		nats.PingInterval(10*time.Second),
		nats.MaxPingsOutstanding(3),
	)
	if err != nil {
		return nil, fmt.Errorf("connect nats server -> %w", err)
	}

	jsContext, err := conn.JetStream(nats.PublishAsyncMaxPending(1024))
	if err != nil {
		return nil, fmt.Errorf("get jetstream context -> %w", err)
	}

	for _, name := range strings.Split(streams, ",") {
		_, err := jsContext.StreamInfo(name)
		if errors.Is(err, nats.ErrStreamNotFound) {
			_, err = jsContext.AddStream(&nats.StreamConfig{
				Name:     name,
				Subjects: []string{fmt.Sprintf("%s.*", name)},
			})
		}
		if err != nil {
			return nil, fmt.Errorf("add stream -> %w", err)
		}
	}

	taskPool, err := ants.NewPoolPreMalloc(1024)
	if err != nil {
		return nil, fmt.Errorf("new task pool -> %w", err)
	}

	consumer := &natsConsumer{conn: conn, jsContext: jsContext, taskPool: taskPool}

	//go consumer.printDetails()

	return consumer, nil
}

func (its *natsConsumer) Subscribe(subj string, queue string, handler IHandler) error {
	var doFunc = func(msg *nats.Msg) {
		err := its.taskPool.Submit(func(msg *nats.Msg) func() {
			return func() {
				err := handler.Process(msg)
				if err != nil {
					log.Error("Error handle message", zap.Error(err))
				}
			}
		}(msg))
		if err != nil {
			log.Error("Error submit task", zap.Error(err))
		}
	}

	t := reflect.TypeOf(handler)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := strings.ReplaceAll(subj, ".", "-") + "-" + strings.ToLower(t.Name())

	if queue == "" {
		_, err := its.jsContext.Subscribe(subj, doFunc, nats.ManualAck(), nats.Durable(name))
		if err != nil {
			return fmt.Errorf("subscribe message -> %w", err)
		}
	} else {
		_, err := its.jsContext.QueueSubscribe(subj, queue, doFunc, nats.ManualAck(), nats.Durable(name))
		if err != nil {
			return fmt.Errorf("queue subscribe message -> %w", err)
		}
	}
	return nil
}

func (its *natsConsumer) printDetails() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Consumer", "Info"})
			for streamInfo := range its.jsContext.Streams() {
				t.AppendRows([]table.Row{{
					"NATS",
					streamInfo.State,
				}})
			}
			t.AppendSeparator()
			t.Render()
		}
	}
}
