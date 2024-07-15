package service

import (
	"context"
	"fmt"

	"github.com/gzericlee/eim/internal/redis"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type GatewayService struct {
	redisManager *redis.Manager
}

func NewGatewayService(redisManager *redis.Manager) *GatewayService {
	return &GatewayService{
		redisManager: redisManager,
	}
}

func (its *GatewayService) RegisterGateway(ctx context.Context, args *rpcmodel.GatewayArgs, reply *rpcmodel.EmptyReply) error {
	err := its.redisManager.RegisterGateway(args.Gateway, args.Expiration)
	if err != nil {
		return fmt.Errorf("save gateway -> %w", err)
	}
	return nil
}

func (its *GatewayService) GetGateways(ctx context.Context, args *rpcmodel.EmptyArgs, reply *rpcmodel.GatewaysReply) error {
	gateways, err := its.redisManager.GetGateways()
	if err != nil {
		return fmt.Errorf("get gateways -> %w", err)
	}

	reply.Gateways = gateways

	return nil
}
