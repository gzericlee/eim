package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type RefresherClient struct {
	*rpcxclient.XClientPool
}

func (its *RefresherClient) RefreshDevicesCache(key string, device *model.Device, action string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Broadcast(ctx, "RefreshDevicesCache", &rpcmodel.RefreshDevicesArgs{Key: key, Device: device, Action: action}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RefreshDevicesCache -> %w", err)
	}
	return nil
}

func (its *RefresherClient) RefreshBizCache(key string, biz *model.Biz, action string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Broadcast(ctx, "RefreshBizCache", &rpcmodel.RefreshBizArgs{Key: key, Biz: biz, Action: action}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RefreshBizsCache -> %w", err)
	}
	return nil
}

func (its *RefresherClient) RefreshBizMembersCache(key string, bizMember *model.BizMember, action string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Broadcast(ctx, "RefreshBizMembersCache", &rpcmodel.RefreshBizMembersArgs{Key: key, BizMember: bizMember, Action: action}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RefreshBizMembersCache -> %w", err)
	}
	return nil
}

func (its *RefresherClient) RefreshTenantCache(key string, tenant *model.Tenant, action string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.Get().Broadcast(ctx, "RefreshTenantCache", &rpcmodel.RefreshTenantArgs{Key: key, Tenant: tenant, Action: action}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RefreshTenantCache -> %w", err)
	}
	return nil
}
