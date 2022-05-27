package seq

import (
	"context"

	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
)

type RpcClient struct {
	pool *client.XClientPool
}

func NewRpcClient(etcdEndpoints []string) (*RpcClient, error) {
	d, err := etcd_client.NewEtcdV3Discovery("/eim_seq", "Id", etcdEndpoints, false, nil)
	if err != nil {
		return nil, err
	}
	pool := client.NewXClientPool(1000, "Id", client.Failover, client.RoundRobin, d, client.DefaultOption)
	return &RpcClient{pool: pool}, nil
}

func (its *RpcClient) Id(userId string) (int64, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "Id", &Request{UserId: userId}, reply)
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}
