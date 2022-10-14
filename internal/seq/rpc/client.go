package rpc

import (
	"context"
	"runtime"

	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
)

type Client struct {
	pool *client.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d, err := etcd_client.NewEtcdV3Discovery(basePath, servicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}
	pool := client.NewXClientPool(runtime.NumCPU(), servicePath, client.Failover, client.RoundRobin, d, client.DefaultOption)
	return &Client{pool: pool}, nil
}

func (its *Client) Number(id string) (int64, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "Number", &Request{Id: id}, reply)
	if err != nil {
		return 0, err
	}
	return reply.Number, nil
}
