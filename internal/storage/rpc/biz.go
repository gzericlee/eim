package rpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"eim/internal/database"
	"eim/internal/model"
	"eim/internal/redis"
	"eim/pkg/cache"
	"eim/util/log"
)

type BizArgs struct {
	Biz *model.Biz
}

type BizReply struct {
	Biz *model.Biz
}

type Biz struct {
	storageCache *cache.Cache
	redisManager *redis.Manager
	database     database.IDatabase
}

func (its *Biz) SaveBiz(ctx context.Context, args *BizArgs, reply *EmptyReply) error {
	err := its.redisManager.SaveBiz(args.Biz)
	if err != nil {
		return fmt.Errorf("save user -> %w", err)
	}

	key := fmt.Sprintf(cacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)

	err = storageRpc.RefreshBizsCache(key, args.Biz, ActionSave)
	if err != nil {
		log.Error("Error refresh biz cache: %v", zap.Error(err))
	}

	return nil
}

func (its *Biz) GetBiz(ctx context.Context, args *BizArgs, reply *BizReply) error {
	key := fmt.Sprintf(cacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)

	if cacheItem, exist := its.storageCache.Get(key); exist {
		reply.Biz = cacheItem.(*model.Biz)
		return nil
	}

	result, err, _ := singleGroup.Do(key, func() (interface{}, error) {
		user, err := its.redisManager.GetBiz(args.Biz.BizId, args.Biz.TenantId)
		if err != nil {
			return nil, fmt.Errorf("get user -> %w", err)
		}
		return user, nil
	})
	if err != nil {
		return fmt.Errorf("biz single group do -> %w", err)
	}

	reply.Biz = result.(*model.Biz)
	its.storageCache.Set(key, reply.Biz)

	return nil
}
