package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type BizClient struct {
	*rpcxclient.XClientPool
}

func (its *BizClient) InsertBiz(biz *model.Biz) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "InsertBiz", &rpcmodel.BizArgs{Biz: biz}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertBiz -> %w", err)
	}

	return nil
}

func (its *BizClient) UpdateBiz(biz *model.Biz) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "UpdateBiz", &rpcmodel.BizArgs{Biz: biz}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call UpdateBiz -> %w", err)
	}

	return nil
}

func (its *BizClient) GetBiz(bizId, tenantId string) (*model.Biz, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.BizReply{}
	err := its.Get().Call(ctx, "GetBiz", &rpcmodel.BizArgs{Biz: &model.Biz{BizId: bizId, TenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetBiz -> %w", err)
	}

	return reply.Biz, nil
}
