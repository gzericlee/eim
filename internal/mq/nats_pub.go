package mq

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"eim/util/log"
)

type natsProducer struct {
	conn *nats.Conn

	successCount int64
	failedCount  int64
	retryCount   int64

	jsContext nats.JetStreamContext
}

func newNatsProducer(endpoints []string) (Producer, error) {
	conn, err := nats.Connect(
		strings.Join(endpoints, ","),
		nats.Name("eim"),
		nats.ReconnectWait(10*time.Second),
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

	producer := &natsProducer{conn: conn, jsContext: jsContext}

	//go producer.printDetails()

	return producer, nil
}

func (its *natsProducer) Publish(subj string, body []byte) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	var err error

	ack, err := its.jsContext.Publish(subj, body)
	if err != nil {
		return fmt.Errorf("publish message -> %w", err)
	}

	log.Debug("Published successfully", zap.String("subject", subj), zap.Any("ack", ack))

	return nil
}

func (its *natsProducer) printDetails() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Producer", "Info"})
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
