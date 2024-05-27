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
	pool *rpcxclient.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, servicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery -> %w", err)
	}
	pool := rpcxclient.NewXClientPool(runtime.NumCPU()*2, servicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d, rpcxclient.DefaultOption)
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
