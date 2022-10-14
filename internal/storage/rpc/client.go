package rpc

import (
	"context"
	"runtime"

	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"

	"eim/internal/types"
)

type Client struct {
	devicePool  *client.XClientPool
	messagePool *client.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d1, err := etcd_client.NewEtcdV3Discovery(basePath, servicePath1, etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}

	d2, err := etcd_client.NewEtcdV3Discovery(basePath, servicePath2, etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		devicePool:  client.NewXClientPool(runtime.NumCPU(), servicePath1, client.Failover, client.RoundRobin, d1, client.DefaultOption),
		messagePool: client.NewXClientPool(runtime.NumCPU(), servicePath2, client.Failover, client.RoundRobin, d2, client.DefaultOption),
	}, nil
}

func (its *Client) SaveDevice(device *types.Device) error {
	err := its.devicePool.Get().Call(context.Background(), "Save", &DeviceRequest{Device: device}, nil)
	return err
}

func (its *Client) SaveMessage(msg *types.Message) error {
	err := its.messagePool.Get().Call(context.Background(), "Save", &MessageRequest{Message: msg}, nil)
	return err
}
