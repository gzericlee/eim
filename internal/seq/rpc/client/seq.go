package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	rpcmodel "github.com/gzericlee/eim/internal/seq/rpc/model"
)

type SeqClient struct {
	*rpcxclient.XClientPool
}

func (its *SeqClient) IncrId(bizId, tenantId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.Reply{}
	err := its.Get().Call(ctx, "IncrId", &rpcmodel.Request{BizId: bizId, TenantId: tenantId}, reply)
	if err != nil {
		return 0, fmt.Errorf("call IncrementId -> %w", err)
	}
	return reply.Number, nil
}

func (its *SeqClient) SnowflakeId() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.Reply{}
	err := its.Get().Call(ctx, "SnowflakeId", &rpcmodel.Request{}, reply)
	if err != nil {
		return 0, fmt.Errorf("call SnowflakeId -> %w", err)
	}

	return reply.Number, nil
}
