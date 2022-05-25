package consumer

import (
	"fmt"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"

	"eim/global"
	"eim/internal/dispatch"
	"eim/internal/gateway/websocket"
	"eim/internal/nsq/api"
	"eim/internal/storage"
	"eim/model"
)

var consumers sync.Map

func InitConsumers(topicChannels map[string][]string, endpoints []string) error {
	config := nsq.NewConfig()
	config.MaxInFlight = 1000

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

	for topic, channels := range topicChannels {
		for _, node := range nodes {
			httpAddr := fmt.Sprintf("%v:%v", node.BroadcastAddress, node.HttpPort)
			err = api.CreateTopic(httpAddr, topic)
			if err != nil {
				return err
			}
			global.Logger.Infof("Created Nsq %v topic %v", httpAddr, topic)
		}
		for _, channel := range channels {
			for _, node := range nodes {
				httpAddr := fmt.Sprintf("%v:%v", node.BroadcastAddress, node.HttpPort)
				err = api.CreateChannel(httpAddr, topic, channel)
				if err != nil {
					return err
				}
				global.Logger.Infof("Created Nsq %v channel %v - %v", httpAddr, topic, channel)
			}

			consumer, err := nsq.NewConsumer(topic, channel, config)
			if err != nil {
				global.Logger.Errorf("Error createing Nsq consumer %v - %v: %s", topic, channel, err)
				return err
			}
			consumer.SetLogger(nil, 0)

			switch topic {
			case model.DeviceStoreTopic:
				{
					consumer.AddConcurrentHandlers(&storage.DeviceHandler{}, config.MaxInFlight)
				}
			case model.MessageStoreTopic:
				{
					consumer.AddConcurrentHandlers(&storage.MessageHandler{}, config.MaxInFlight)
				}
			case model.MessageDispatchTopic:
				{
					consumer.AddConcurrentHandlers(&dispatch.MessageHandler{}, config.MaxInFlight)
				}
			case model.MessageSendTopic:
				{
					consumer.AddConcurrentHandlers(&websocket.SendHandler{}, config.MaxInFlight)
				}
			}

			err = consumer.ConnectToNSQLookupds(endpoints)
			if err != nil {
				global.Logger.Errorf("Error connecting to Nsq %v : %s", endpoints, err)
				time.Sleep(time.Second)
			}

			consumers.Store(fmt.Sprintf("%v - %v", topic, channel), consumer)
		}
	}

	//printConsumersDetail()

	return nil
}
