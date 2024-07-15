package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type TenantClient struct {
	*rpcxclient.XClientPool
}

func (its *TenantClient) InsertTenant(tenant *model.Tenant) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "InsertTenant", &rpcmodel.TenantArgs{Tenant: tenant}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertTenant -> %w", err)
	}
	return nil
}

func (its *TenantClient) UpdateTenant(tenant *model.Tenant) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "UpdateTenant", &rpcmodel.TenantArgs{Tenant: tenant}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call UpdateTenant -> %w", err)
	}

	return nil
}

func (its *TenantClient) GetTenant(tenantId string) (*model.Tenant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.TenantReply{}
	err := its.Get().Call(ctx, "GetTenant", &rpcmodel.TenantArgs{Tenant: &model.Tenant{TenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetTenant -> %w", err)
	}

	return reply.Tenant, nil
}
