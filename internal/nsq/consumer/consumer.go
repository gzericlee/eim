package consumer

import (
	"fmt"
	"sync"

	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

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
			global.Logger.Warn("Error getting Nsq nodes", zap.String("endpoint", endpoint), zap.Error(err))
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
			global.Logger.Info("Created Nsq topic", zap.String("endpoint", httpAddr), zap.String("topic", topic))
		}
		for _, channel := range channels {
			for _, node := range nodes {
				httpAddr := fmt.Sprintf("%v:%v", node.BroadcastAddress, node.HttpPort)
				err = api.CreateChannel(httpAddr, topic, channel)
				if err != nil {
					return err
				}
				global.Logger.Info("Created Nsq channel", zap.String("endpoint", httpAddr), zap.String("topic", topic), zap.String("channel", channel))
			}

			consumer, err := nsq.NewConsumer(topic, channel, config)
			if err != nil {
				global.Logger.Error("Error creating Nsq consumer", zap.String("topic", topic), zap.String("channel", channel), zap.Error(err))
				return err
			}
			consumer.SetLogger(nil, 0)

			storageRpc, err := storage.NewRpcClient(global.SystemConfig.Etcd.Endpoints.Value())
			if err != nil {
				return err
			}

			switch topic {
			case model.DeviceStoreTopic:
				{
					consumer.AddConcurrentHandlers(&dispatch.DeviceHandler{StorageRpc: storageRpc}, config.MaxInFlight)
				}
			case model.MessageDispatchTopic:
				{
					consumer.AddConcurrentHandlers(&dispatch.MessageHandler{StorageRpc: storageRpc}, config.MaxInFlight)
				}
			case model.MessageSendTopic:
				{
					consumer.AddConcurrentHandlers(&websocket.SendHandler{}, config.MaxInFlight)
				}
			}

			err = consumer.ConnectToNSQLookupds(endpoints)
			if err != nil {
				return err
			}

			consumers.Store(fmt.Sprintf("%v - %v", topic, channel), consumer)
		}
	}

	//printConsumersDetail()

	return nil
}
