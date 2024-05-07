package rpc

import (
	"context"
	"runtime"

	rpcx_etcd_client "github.com/rpcxio/rpcx-etcd/client"
	rpcx_client "github.com/smallnest/rpcx/client"
)

type Client struct {
	pool *rpcx_client.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d, err := rpcx_etcd_client.NewEtcdV3Discovery(basePath, servicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}
	pool := rpcx_client.NewXClientPool(runtime.NumCPU(), servicePath, rpcx_client.Failover, rpcx_client.ConsistentHash, d, rpcx_client.DefaultOption)
	return &Client{pool: pool}, nil
}

func (its *Client) IncrementId(bizId string) (int64, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "IncrementId", &Request{BizId: bizId}, reply)
	if err != nil {
		return 0, err
	}
	return reply.Number, nil
}

func (its *Client) SnowflakeId() (int64, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "SnowflakeId", &Request{}, reply)
	if err != nil {
		return 0, err
	}
	return reply.Number, nil
}
