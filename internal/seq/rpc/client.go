package rpc

import (
	"context"
	"fmt"
	"runtime"

	rpcxetcdclient "github.com/rpcxio/rpcx-etcd/client"
	rpcxclient "github.com/smallnest/rpcx/client"
)

type Client struct {
	pool *rpcxclient.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, servicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery -> %w", err)
	}
	pool := rpcxclient.NewXClientPool(runtime.NumCPU()*2, servicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d, rpcxclient.DefaultOption)
	return &Client{pool: pool}, nil
}

func (its *Client) IncrementId(bizId string) (int64, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "IncrementId", &Request{BizId: bizId}, reply)
	if err != nil {
		return 0, fmt.Errorf("call increment id -> %w", err)
	}
	return reply.Number, nil
}

func (its *Client) SnowflakeId() (int64, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "SnowflakeId", &Request{}, reply)
	if err != nil {
		return 0, fmt.Errorf("call snowflake id -> %w", err)
	}
	return reply.Number, nil
}
