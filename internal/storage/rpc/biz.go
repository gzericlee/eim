package rpc

import (
	"context"
	"fmt"
	"time"

	"eim/internal/database"
	"eim/internal/model"
	"eim/internal/redis"
	"eim/pkg/cache"
	"eim/pkg/cache/notify"
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
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.redisManager.SaveBiz(args.Biz)
	if err != nil {
		return fmt.Errorf("save user -> %w", err)
	}

	key := fmt.Sprintf(cacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)
	err = notify.Del(bizCachePool, key)
	if err != nil {
		return fmt.Errorf("del user(%s) cache -> %w", key, err)
	}

	return nil
}

func (its *Biz) GetBiz(ctx context.Context, args *BizArgs, reply *BizReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	key := fmt.Sprintf(cacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)

	if cacheItem, exist := its.storageCache.Get(key); exist {
		reply.Biz = cacheItem.(*model.Biz)
		return nil
	}

	result, err, _ := group.Do(key, func() (interface{}, error) {
		user, err := its.redisManager.GetBiz(args.Biz.BizId, args.Biz.TenantId)
		if err != nil {
			return nil, fmt.Errorf("get user -> %w", err)
		}
		return user, nil
	})
	if err != nil {
		return fmt.Errorf("group do -> %w", err)
	}

	reply.Biz = result.(*model.Biz)
	its.storageCache.Put(key, reply.Biz)

	return nil
}
