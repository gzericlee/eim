package rpc

import (
	"context"
	"fmt"
	"time"

	"eim/internal/model"
	"eim/internal/redis"
	"eim/util/log"
)

type RegisterGatewayArgs struct {
	Gateway    *model.Gateway
	Expiration time.Duration
}

type GetGatewaysReply struct {
	Gateways []*model.Gateway
}

type Gateway struct {
	redisManager *redis.Manager
}

func (its *Gateway) RegisterGateway(ctx context.Context, args *RegisterGatewayArgs, reply *EmptyReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.redisManager.RegisterGateway(args.Gateway, args.Expiration)
	if err != nil {
		return fmt.Errorf("save gateway -> %w", err)
	}

	return nil
}

func (its *Gateway) GetGateways(ctx context.Context, args *EmptyArgs, reply *GetGatewaysReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	gateways, err := its.redisManager.GetGateways()
	if err != nil {
		return fmt.Errorf("get gateways -> %w", err)
	}

	reply.Gateways = gateways

	return nil
}
