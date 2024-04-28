package rpc

import (
	"context"
	"runtime"

	rpcx_etcd_client "github.com/rpcxio/rpcx-etcd/client"
	rpcx_client "github.com/smallnest/rpcx/client"

	"eim/internal/model"
)

type Client struct {
	devicePool  *rpcx_client.XClientPool
	messagePool *rpcx_client.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d1, err := rpcx_etcd_client.NewEtcdV3Discovery(basePath, servicePath1, etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}

	d2, err := rpcx_etcd_client.NewEtcdV3Discovery(basePath, servicePath2, etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		devicePool:  rpcx_client.NewXClientPool(runtime.NumCPU(), servicePath1, rpcx_client.Failover, rpcx_client.RoundRobin, d1, rpcx_client.DefaultOption),
		messagePool: rpcx_client.NewXClientPool(runtime.NumCPU(), servicePath2, rpcx_client.Failover, rpcx_client.RoundRobin, d2, rpcx_client.DefaultOption),
	}, nil
}

func (its *Client) SaveDevice(device *model.Device) error {
	err := its.devicePool.Get().Call(context.Background(), "Save", &DeviceRequest{Device: device}, nil)
	return err
}

func (its *Client) SaveMessage(msg *model.Message) error {
	err := its.messagePool.Get().Call(context.Background(), "Save", &MessageRequest{Message: msg}, nil)
	return err
}
