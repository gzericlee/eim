package rpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"eim/internal/database"
	"eim/internal/model"
	"eim/pkg/cache"
	"eim/pkg/log"
)

type TenantArgs struct {
	Tenant *model.Tenant
}

type TenantReply struct {
	Tenant *model.Tenant
}

type Tenant struct {
	storageCache *cache.Cache[string, *model.Tenant]
	database     database.IDatabase
}

func (its *Tenant) InsertTenant(ctx context.Context, args *TenantArgs, reply *EmptyReply) error {
	err := its.database.InsertTenant(args.Tenant)
	if err != nil {
		return fmt.Errorf("insert tenant -> %w", err)
	}

	key := fmt.Sprintf(tenantCacheKeyFormat, tenantCachePool, args.Tenant.TenantId)

	err = storageRpc.RefreshTenantCache(key, args.Tenant, ActionSave)
	if err != nil {
		log.Error("Error refresh tenant cache: %v", zap.Error(err))
	}

	return nil
}

func (its *Tenant) UpdateTenant(ctx context.Context, args *TenantArgs, reply *EmptyReply) error {
	err := its.database.UpdateTenant(args.Tenant)
	if err != nil {
		return fmt.Errorf("update tenant -> %w", err)
	}

	key := fmt.Sprintf(tenantCacheKeyFormat, tenantCachePool, args.Tenant.TenantId)

	err = storageRpc.RefreshTenantCache(key, args.Tenant, ActionSave)
	if err != nil {
		log.Error("Error refresh tenant cache: %v", zap.Error(err))
	}

	return nil
}

func (its *Tenant) GetTenant(ctx context.Context, args *TenantArgs, reply *TenantReply) error {
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
