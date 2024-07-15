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

type BizService struct {
	storageCache *cache.Cache[string, *model.Biz]
	database     database.IDatabase
	refreshRpc   *rpcclient.RefresherClient
}

func NewBizService(storageCache *cache.Cache[string, *model.Biz], database database.IDatabase, refreshRpc *rpcclient.RefresherClient) *BizService {
	return &BizService{
		storageCache: storageCache,
		database:     database,
		refreshRpc:   refreshRpc,
	}
}

func (its *BizService) InsertBiz(ctx context.Context, args *rpcmodel.BizArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.InsertBiz(args.Biz)
	if err != nil {
		return fmt.Errorf("insert biz -> %w", err)
	}

	key := fmt.Sprintf(bizCacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)

	err = its.refreshRpc.RefreshBizCache(key, args.Biz, rpcmodel.ActionInsert)
	if err != nil {
		log.Error("Error refresh biz cache: %v", zap.Error(err))
	}

	return nil
}

func (its *BizService) UpdateBiz(ctx context.Context, args *rpcmodel.BizArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.UpdateBiz(args.Biz)
	if err != nil {
		return fmt.Errorf("update biz -> %w", err)
	}

	key := fmt.Sprintf(bizCacheKeyFormat, bizCachePool, args.Biz.BizId, args.Biz.TenantId)

	err = its.refreshRpc.RefreshBizCache(key, args.Biz, rpcmodel.ActionUpdate)
	if err != nil {
		log.Error("Error refresh biz cache: %v", zap.Error(err))
	}

	return nil
}

func (its *BizService) GetBiz(ctx context.Context, args *rpcmodel.BizArgs, reply *rpcmodel.BizReply) error {
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
