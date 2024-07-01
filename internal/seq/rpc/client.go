package rpc

import (
	"context"
	"fmt"
	"runtime"
	"time"

	rpcxetcdclient "github.com/rpcxio/rpcx-etcd/client"
	rpcxclient "github.com/smallnest/rpcx/client"
)

type Client struct {
	pool *rpcxclient.XClientPool
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, servicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for seq service -> %w", err)
	}
	pool := rpcxclient.NewXClientPool(runtime.NumCPU()*2, servicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d, rpcxclient.DefaultOption)
	return &Client{pool: pool}, nil
}

func (its *Client) IncrId(bizId, tenantId string) (int64, error) {
	reply := &Reply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.pool.Get().Call(ctx, "IncrId", &Request{BizId: bizId, TenantId: tenantId}, reply)
	if err != nil {
		return 0, fmt.Errorf("call IncrementId -> %w", err)
	}
	return reply.Number, nil
}

func (its *Client) SnowflakeId() (int64, error) {
	reply := &Reply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.pool.Get().Call(ctx, "SnowflakeId", &Request{}, reply)
	if err != nil {
		return 0, fmt.Errorf("call SnowflakeId -> %w", err)
	}
	return reply.Number, nil
}
