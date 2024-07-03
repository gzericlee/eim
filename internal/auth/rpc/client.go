package rpc

import (
	"context"
	"fmt"
	"time"

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
		return nil, fmt.Errorf("new etcd v3 discovery for auth service -> %w", err)
	}
	pool := rpcxclient.NewXClientPool(100, servicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d, rpcxclient.DefaultOption)
	return &Client{pool: pool}, nil
}

func (its *Client) CheckToken(token string) (*model.Biz, error) {
	reply := &Reply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.pool.Get().Call(ctx, "CheckToken", &Request{Token: token}, reply)
	if err != nil {
		return nil, fmt.Errorf("call CheckToken -> %w", err)
	}
	return reply.Biz, nil
}
