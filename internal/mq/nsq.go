package mq

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/internal/mq/api"
	"eim/pkg/log"
)

type nsqProducer struct {
	pool     *producerPool
	doneChan chan *nsq.ProducerTransaction
}

type nsqConsumer struct {
	nodes     []*api.Node
	endpoints []string
	locker    sync.RWMutex
	consumers map[string]*nsq.Consumer
}

func (its *nsqProducer) Publish(topic string, body []byte) error {
	err := its.pool.publishAsync(topic, body, its.doneChan)
	if err != nil {
		return err
	}
	log.Debug("Published successfully", zap.String("topic", topic), zap.ByteString("body", body))
	return nil
}

func (its *nsqConsumer) Subscribe(topic string, channel string, handler nsq.Handler) error {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = runtime.NumCPU() * 100   // 根据你的CPU核心数调整
	cfg.MsgTimeout = time.Second * 60          // 设置为你的消息处理的最大预期时间
	cfg.HeartbeatInterval = time.Second * 30   // 通常设置为10到60秒之间
	cfg.MaxRequeueDelay = time.Second * 15     // 根据你的消息处理失败后的预期重试时间来调整
	cfg.DefaultRequeueDelay = time.Second * 3  // 根据你的消息处理失败后的预期重试时间来调整
	cfg.MaxBackoffDuration = time.Second * 120 // 根据你的应用程序的特性来调整
	cfg.LookupdPollInterval = time.Second * 60 // 通常设置为你的NSQLookupd更新频率的两倍
	cfg.LookupdPollJitter = 0.3                // 通常设置为0.1到0.3之间
	cfg.MaxAttempts = 5                        // 根据你的消息处理失败后的预期重试次数来调整
	cfg.LowRdyIdleTimeout = time.Second * 10   // 根据你的消息处理的最大预期时间来调整

	for _, node := range its.nodes {
		httpAddr := fmt.Sprintf("%v:%v", node.BroadcastAddress, node.HttpPort)
		err := api.CreateTopic(httpAddr, topic)
		if err != nil {
			return err
		}
		log.Info("Created nsq topic successfully", zap.String("endpoint", httpAddr), zap.String("topic", topic))
		err = api.CreateChannel(httpAddr, topic, channel)
		if err != nil {
			return err
		}
		log.Info("Created nsq channel successfully", zap.String("endpoint", httpAddr), zap.String("topic", topic), zap.String("channel", channel))
	}

	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Error("Error creating nsq consumer", zap.String("topic", topic), zap.String("channel", channel), zap.Error(err))
		return err
	}
	consumer.SetLogger(nil, 0)

	consumer.AddHandler(handler)

	err = consumer.ConnectToNSQLookupds(its.endpoints)
	if err != nil {
		return err
	}

	its.locker.Lock()
	its.consumers[fmt.Sprintf("%v - %v", topic, channel)] = consumer
	its.locker.Unlock()

	return nil
}

func newNsqProducer(nsqEndpoints []string) (Producer, error) {
	var err error
	var nodes []*api.Node

	for _, endpoint := range nsqEndpoints {
		nodes, err = api.GetNodes(endpoint)
		if err != nil {
			log.Warn("Error getting nsq nodes", zap.String("endpoint", endpoint), zap.Error(err))
			continue
		}
		break
	}
	if err != nil {
		return nil, fmt.Errorf("cannot find nsq %v nodes", nsqEndpoints)
	}

	var adders []string
	for _, node := range nodes {
		tcpAddr := fmt.Sprintf("%v:%v", node.BroadcastAddress, node.TcpPort)
		adders = append(adders, tcpAddr)
	}

	config := nsq.NewConfig()

	pool, err := newProducerPool(adders, config, 32)
	if err != nil {
		return nil, err
	}

	producer := &nsqProducer{
		pool:     pool,
		doneChan: make(chan *nsq.ProducerTransaction, 1),
	}

	go func() {
		for done := range producer.doneChan {
			if done.Error != nil {
				log.Error("Error publishing to nsq", zap.Error(err))
			}
		}
	}()

	go producer.printDetails()

	return producer, err
}

func newNsqConsumer(nsqEndpoints []string) (Consumer, error) {
	var err error
	var nodes []*api.Node
	for _, endpoint := range nsqEndpoints {
		nodes, err = api.GetNodes(endpoint)
		if err != nil {
			log.Warn("Error getting nsq nodes", zap.String("endpoint", endpoint), zap.Error(err))
			continue
		}
		break
	}
	if err != nil {
		return nil, fmt.Errorf("cannot find nsq %v nodes", nsqEndpoints)
	}

	consumer := &nsqConsumer{
		nodes:     nodes,
		endpoints: nsqEndpoints,
		consumers: make(map[string]*nsq.Consumer),
	}

	go consumer.printDetails()

	return consumer, nil
}

func (its *nsqProducer) printDetails() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Producer", "Successful", "Failed", "Retry"})
			t.AppendRows([]table.Row{{
				"NSQ",
				its.pool.successCount,
				its.pool.failedCount,
				its.pool.retryCount,
			}})
			t.AppendSeparator()
			t.Render()
		}
	}
}

func (its *nsqConsumer) printDetails() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Consumer", "Connections", "Received", "Finished", "Requeued"})
			for name, consumer := range its.consumers {
				stats := consumer.Stats()
				t.AppendRows([]table.Row{{
					name,
					stats.Connections,
					stats.MessagesReceived,
					stats.MessagesFinished,
					stats.MessagesRequeued,
				}})
			}
			t.AppendSeparator()
			t.Render()
		}
	}
}
