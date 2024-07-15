package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type BizMemberClient struct {
	*rpcxclient.XClientPool
}

func (its *BizMemberClient) AddBizMember(member *model.BizMember) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "AddBizMember", &rpcmodel.BizMemberArgs{BizMember: member}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call AddBizMember -> %w", err)
	}

	return nil
}

func (its *BizMemberClient) RemoveBizMember(member *model.BizMember) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "RemoveBizMember", &rpcmodel.BizMemberArgs{BizMember: member}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RemoveBizMember -> %w", err)
	}

	return nil
}

func (its *BizMemberClient) GetBizMembers(bizId, tenantId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.BizMembersReply{}
	err := its.Get().Call(ctx, "GetBizMembers", &rpcmodel.BizMemberArgs{BizMember: &model.BizMember{BizId: bizId, BizTenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetBizMembers -> %w", err)
	}

	return reply.Members, nil
}
