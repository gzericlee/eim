package rpc

import (
	"context"
	"fmt"
	"runtime"

	rpcxetcdclient "github.com/rpcxio/rpcx-etcd/client"
	rpcxclient "github.com/smallnest/rpcx/client"

	"eim/internal/model"
)

type Client struct {
	devicePool  *rpcxclient.XClientPool
	messagePool *rpcxclient.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d1, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, servicePath1, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for device -> %w", err)
	}

	d2, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, servicePath2, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for message -> %w", err)
	}

	return &Client{
		devicePool:  rpcxclient.NewXClientPool(runtime.NumCPU(), servicePath1, rpcxclient.Failover, rpcxclient.RoundRobin, d1, rpcxclient.DefaultOption),
		messagePool: rpcxclient.NewXClientPool(runtime.NumCPU(), servicePath2, rpcxclient.Failover, rpcxclient.RoundRobin, d2, rpcxclient.DefaultOption),
	}, nil
}

func (its *Client) SaveDevice(device *model.Device) error {
	err := its.devicePool.Get().Call(context.Background(), "Save", &DeviceRequest{Device: device}, nil)
	if err != nil {
		return fmt.Errorf("call save device -> %w", err)
	}
	return nil
}

func (its *Client) SaveMessage(msg *model.Message) error {
	err := its.messagePool.Get().Call(context.Background(), "Save", &MessageRequest{Message: msg}, nil)
	if err != nil {
		return fmt.Errorf("call save message -> %w", err)
	}
	return nil
}
