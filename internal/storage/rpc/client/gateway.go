package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type GatewayClient struct {
	*rpcxclient.XClientPool
}

func (its *GatewayClient) RegisterGateway(gateway *model.Gateway, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "RegisterGateway", &rpcmodel.GatewayArgs{Gateway: gateway, Expiration: expiration}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RegisterGateway -> %w", err)
	}
	return nil
}

func (its *GatewayClient) GetGateways() ([]*model.Gateway, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.GatewaysReply{}
	err := its.Get().Call(ctx, "GetGateways", &rpcmodel.EmptyArgs{}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetGateways -> %w", err)
	}

	return reply.Gateways, nil
}
