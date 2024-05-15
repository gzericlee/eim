package rpc

import (
	"context"
	"fmt"
	"runtime"

	rpcx_etcd_client "github.com/rpcxio/rpcx-etcd/client"
	rpcx_client "github.com/smallnest/rpcx/client"

	"eim/internal/model"
)

type Client struct {
	pool *rpcx_client.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d, err := rpcx_etcd_client.NewEtcdV3Discovery(basePath, servicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery -> %w", err)
	}
	pool := rpcx_client.NewXClientPool(runtime.NumCPU(), servicePath, rpcx_client.Failover, rpcx_client.RoundRobin, d, rpcx_client.DefaultOption)
	return &Client{pool: pool}, nil
}

func (its *Client) CheckToken(token string) (*model.User, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "CheckToken", &Request{Token: token}, reply)
	if err != nil {
		return nil, fmt.Errorf("call check token -> %w", err)
	}
	return reply.User, nil
}
