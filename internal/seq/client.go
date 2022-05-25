package seq

import (
	"context"

	"github.com/smallnest/rpcx/client"
)

type RpcClient struct {
	endpoint string
	pool     *client.XClientPool
}

func NewRpcClient(endpoint string) (*RpcClient, error) {
	d, err := client.NewPeer2PeerDiscovery("tcp@"+endpoint, "")
	if err != nil {
		return nil, err
	}
	pool := client.NewXClientPool(1000, "seq.id.service", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	return &RpcClient{endpoint: endpoint, pool: pool}, nil
}

func (its *RpcClient) ID(userId string) (int64, error) {
	reply := &Reply{}
	err := its.pool.Get().Call(context.Background(), "id", &Request{UserId: userId}, reply)
	if err != nil {
		return 0, err
	}
	return reply.Id, nil
}
