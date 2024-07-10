package rpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"eim/internal/database"
	"eim/internal/model"
	"eim/internal/redis"
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
	redisManager *redis.Manager
	database     database.IDatabase
}

func (its *Tenant) SaveTenant(ctx context.Context, args *TenantArgs, reply *EmptyReply) error {
	err := its.redisManager.SaveTenant(args.Tenant)
	if err != nil {
		return fmt.Errorf("save tenant -> %w", err)
	}

	key := fmt.Sprintf(tenantCacheKeyFormat, tenantCachePool, args.Tenant.TenantId)

	err = storageRpc.RefreshTenantCache(key, args.Tenant, ActionSave)
	if err != nil {
		log.Error("Error refresh biz cache: %v", zap.Error(err))
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
		user, err := its.redisManager.GetTenant(args.Tenant.TenantId)
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
