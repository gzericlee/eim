package rpc

import (
	"context"
	"fmt"
	"time"

	"eim/internal/model"
	"eim/internal/redis"
)

type GatewayArgs struct {
	Gateway    *model.Gateway
	Expiration time.Duration
}

type GatewaysReply struct {
	Gateways []*model.Gateway
}

type Gateway struct {
	redisManager *redis.Manager
}

func (its *Gateway) RegisterGateway(ctx context.Context, args *GatewayArgs, reply *EmptyReply) error {
	err := its.redisManager.RegisterGateway(args.Gateway, args.Expiration)
	if err != nil {
		return fmt.Errorf("save gateway -> %w", err)
	}
	return nil
}

func (its *Gateway) GetGateways(ctx context.Context, args *EmptyArgs, reply *GatewaysReply) error {
	gateways, err := its.redisManager.GetGateways()
	if err != nil {
		return fmt.Errorf("get gateways -> %w", err)
	}

	reply.Gateways = gateways

	return nil
}
