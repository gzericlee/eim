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

type BizArgs struct {
	Biz *model.Biz
}

type BizReply struct {
	Biz *model.Biz
}

type Biz struct {
	storageCache *cache.Cache[string, *model.Biz]
	database     database.IDatabase
}

func (its *Biz) InsertBiz(ctx context.Context, args *BizArgs, reply *EmptyReply) error {
	err := its.database.InsertBiz(args.Biz)
	if err != nil {
		return fmt.Errorf("insert biz -> %w", err)
	}

	key := fmt.Sprintf(bizCacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)

	err = storageRpc.RefreshBizCache(key, args.Biz, ActionSave)
	if err != nil {
		log.Error("Error refresh biz cache: %v", zap.Error(err))
	}

	return nil
}

func (its *Biz) UpdateBiz(ctx context.Context, args *BizArgs, reply *EmptyReply) error {
	err := its.database.UpdateBiz(args.Biz)
	if err != nil {
		return fmt.Errorf("update biz -> %w", err)
	}

	key := fmt.Sprintf(bizCacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)

	err = storageRpc.RefreshBizCache(key, args.Biz, ActionSave)
	if err != nil {
		log.Error("Error refresh biz cache: %v", zap.Error(err))
	}

	return nil
}

func (its *Biz) GetBiz(ctx context.Context, args *BizArgs, reply *BizReply) error {
	key := fmt.Sprintf(bizCacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)

	if biz, exist := its.storageCache.Get(key); exist {
		reply.Biz = biz
		return nil
	}

	result, err, _ := singleGroup.Do(key, func() (interface{}, error) {
		biz, err := its.database.GetBiz(args.Biz.BizId, args.Biz.TenantId)
		if err != nil {
			return nil, fmt.Errorf("get biz -> %w", err)
		}
		return biz, nil
	})
	if err != nil {
		return fmt.Errorf("single group do -> %w", err)
	}

	reply.Biz = result.(*model.Biz)
	its.storageCache.Set(key, reply.Biz)

	return nil
}
