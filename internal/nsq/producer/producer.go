package producer

import (
	"fmt"

	"github.com/nsqio/go-nsq"

	"eim/global"
	"eim/internal/nsq/api"
)

var pool *Pool
var doneChan = make(chan *nsq.ProducerTransaction, 1)

func Publish(topic string, body []byte) error {
	err := pool.Publish(topic, body)
	if err != nil {
		return err
	}
	global.Logger.Debugf("Published to %v successful，body: %v", topic, string(body))
	return nil
}

func PublishAsync(topic string, body []byte) error {
	err := pool.PublishAsync(topic, body, doneChan)
	if err != nil {
		return err
	}
	global.Logger.Debugf("Published to %v successful，body: %v", topic, string(body))
	return nil
}

func InitProducers(endpoints []string) error {
	var err error
	var nodes []*api.Node

	for _, endpoint := range endpoints {
		nodes, err = api.GetNodes(endpoint)
		if err != nil {
			global.Logger.Warnf("Error geting Nsq %v nodes: %v", endpoint, err)
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("cannot find Nsq %v nodes", endpoints)
	}

	var addrs []string
	for _, node := range nodes {
		tcpAddr := fmt.Sprintf("%v:%v", node.BroadcastAddress, node.TcpPort)
		addrs = append(addrs, tcpAddr)
	}

	config := nsq.NewConfig()
	config.MaxInFlight = 1000

	pool, err = NewPool(addrs, config)
	pool.MaxAttempts = len(addrs)

	go func() {
		for done := range doneChan {
			if done.Error != nil {
				global.Logger.Errorf("Error publishing to Nsq: %v", err)
			}
		}
	}()

	return err
}
