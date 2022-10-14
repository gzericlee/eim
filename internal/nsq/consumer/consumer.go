package consumer

import (
	"fmt"
	"sync"

	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/dispatch"
	"eim/internal/gateway/websocket"
	"eim/internal/nsq/api"
	storage_rpc "eim/internal/storage/rpc"
	"eim/internal/types"
	"eim/pkg/log"
)

var consumers sync.Map

func InitConsumers(topicChannels map[string][]string, endpoints []string) error {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000

	var err error
	var nodes []*api.Node
	for _, endpoint := range endpoints {
		nodes, err = api.GetNodes(endpoint)
		if err != nil {
			log.Warn("Error getting Nsq nodes", zap.String("endpoint", endpoint), zap.Error(err))
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
			log.Info("Created Nsq topic", zap.String("endpoint", httpAddr), zap.String("topic", topic))
		}
		for _, channel := range channels {
			for _, node := range nodes {
				httpAddr := fmt.Sprintf("%v:%v", node.BroadcastAddress, node.HttpPort)
				err = api.CreateChannel(httpAddr, topic, channel)
				if err != nil {
					return err
				}
				log.Info("Created Nsq channel", zap.String("endpoint", httpAddr), zap.String("topic", topic), zap.String("channel", channel))
			}

			consumer, err := nsq.NewConsumer(topic, channel, cfg)
			if err != nil {
				log.Error("Error creating Nsq consumer", zap.String("topic", topic), zap.String("channel", channel), zap.Error(err))
				return err
			}
			consumer.SetLogger(nil, 0)

			storageRpc, err := storage_rpc.NewClient(config.SystemConfig.Etcd.Endpoints.Value())
			if err != nil {
				return err
			}

			switch topic {
			case types.DeviceStoreTopic:
				{
					consumer.AddConcurrentHandlers(&dispatch.DeviceHandler{StorageRpc: storageRpc}, cfg.MaxInFlight)
				}
			case types.MessageUserDispatchTopic:
				{
					consumer.AddConcurrentHandlers(&dispatch.UserMessageHandler{StorageRpc: storageRpc}, cfg.MaxInFlight)
				}
			case types.MessageGroupDispatchTopic:
				{
					consumer.AddConcurrentHandlers(&dispatch.GroupMessageHandler{StorageRpc: storageRpc}, cfg.MaxInFlight)
				}
			case types.MessageSendTopic:
				{
					consumer.AddConcurrentHandlers(&websocket.SendHandler{}, cfg.MaxInFlight)
				}
			}

			err = consumer.ConnectToNSQLookupds(endpoints)
			if err != nil {
				return err
			}

			consumers.Store(fmt.Sprintf("%v - %v", topic, channel), consumer)
		}
	}

	printConsumersDetail()

	return nil
}
