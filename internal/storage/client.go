package storage

import (
	"context"

	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"

	"eim/model"
	"eim/proto/pb"
)

type RpcClient struct {
	devicePool  *client.XClientPool
	messagePool *client.XClientPool
}

func NewRpcClient(etcdEndpoints []string) (*RpcClient, error) {
	d1, err := etcd_client.NewEtcdV3Discovery("/eim_storage", "Device", etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}

	d2, err := etcd_client.NewEtcdV3Discovery("/eim_storage", "Message", etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		devicePool:  client.NewXClientPool(1000, "Device", client.Failover, client.RoundRobin, d1, client.DefaultOption),
		messagePool: client.NewXClientPool(1000, "Message", client.Failover, client.RoundRobin, d2, client.DefaultOption),
	}, nil
}

func (its *RpcClient) SaveDevice(device *model.Device) error {
	err := its.devicePool.Get().Call(context.Background(), "Save", &DeviceRequest{Device: device}, nil)
	return err
}

func (its *RpcClient) SaveMessage(msg *pb.Message) error {
	err := its.messagePool.Get().Call(context.Background(), "Save", &MessageRequest{Message: msg}, nil)
	return err
}
