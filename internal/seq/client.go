package seq

import (
	"context"
	"runtime"

	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
)

type RpcClient struct {
	pool *client.XClientPool
}

func NewRpcClient(etcdEndpoints []string) (*RpcClient, error) {
	d, err := etcd_client.NewEtcdV3Discovery(basePath, servicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}
	pool := client.NewXClientPool(runtime.NumCPU(), servicePath, client.Failover, client.RoundRobin, d, client.DefaultOption)
	return &RpcClient{pool: pool}, nil
}

func (its *RpcClient) Number(id string) (int64, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "Number", &Request{Id: id}, reply)
	if err != nil {
		return 0, err
	}
	return reply.Number, nil
}
