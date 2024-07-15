package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/gzericlee/eim/internal/database"
	"github.com/gzericlee/eim/internal/model"
	rpcclient "github.com/gzericlee/eim/internal/storage/rpc/client"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
	"github.com/gzericlee/eim/pkg/cache"
	"github.com/gzericlee/eim/pkg/log"
)

type TenantService struct {
	storageCache *cache.Cache[string, *model.Tenant]
	database     database.IDatabase
	refreshRpc   *rpcclient.RefresherClient
}

func NewTenantService(storageCache *cache.Cache[string, *model.Tenant], database database.IDatabase, refreshRpc *rpcclient.RefresherClient) *TenantService {
	return &TenantService{
		storageCache: storageCache,
		database:     database,
		refreshRpc:   refreshRpc,
	}
}

func (its *TenantService) InsertTenant(ctx context.Context, args *rpcmodel.TenantArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.InsertTenant(args.Tenant)
	if err != nil {
		return fmt.Errorf("insert tenant -> %w", err)
	}

	key := fmt.Sprintf(tenantCacheKeyFormat, tenantCachePool, args.Tenant.TenantId)

	err = its.refreshRpc.RefreshTenantCache(key, args.Tenant, rpcmodel.ActionInsert)
	if err != nil {
		log.Error("Error refresh tenant cache: %v", zap.Error(err))
	}

	return nil
}

func (its *TenantService) UpdateTenant(ctx context.Context, args *rpcmodel.TenantArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.UpdateTenant(args.Tenant)
	if err != nil {
		return fmt.Errorf("update tenant -> %w", err)
	}

	key := fmt.Sprintf(tenantCacheKeyFormat, tenantCachePool, args.Tenant.TenantId)

	err = its.refreshRpc.RefreshTenantCache(key, args.Tenant, rpcmodel.ActionUpdate)
	if err != nil {
		log.Error("Error refresh tenant cache: %v", zap.Error(err))
	}

	return nil
}

func (its *TenantService) GetTenant(ctx context.Context, args *rpcmodel.TenantArgs, reply *rpcmodel.TenantReply) error {
	key := fmt.Sprintf(tenantCacheKeyFormat, bizCachePool, args.Tenant.TenantId)

	if tenant, exist := its.storageCache.Get(key); exist {
		reply.Tenant = tenant
		return nil
	}

	result, err, _ := singleGroup.Do(key, func() (interface{}, error) {
		user, err := its.database.GetTenant(args.Tenant.TenantId)
		if err != nil {
			return nil, fmt.Errorf("get tenant -> %w", err)
		}
		return user, nil
	})
	if err != nil {
		return fmt.Errorf("single group do -> %w", err)
	}

	reply.Tenant = result.(*model.Tenant)
	its.storageCache.Set(key, reply.Tenant)

	return nil
}
