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

type BizMemberService struct {
	storageCache *cache.Cache[string, []*model.BizMember]
	database     database.IDatabase
	refreshRpc   *rpcclient.RefresherClient
}

func NewBizMemberService(storageCache *cache.Cache[string, []*model.BizMember], database database.IDatabase, refreshRpc *rpcclient.RefresherClient) *BizMemberService {
	return &BizMemberService{
		storageCache: storageCache,
		database:     database,
		refreshRpc:   refreshRpc,
	}
}

func (its *BizMemberService) AddBizMember(ctx context.Context, args *rpcmodel.BizMemberArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.InsertBizMember(args.BizMember)
	if err != nil {
		return fmt.Errorf("add biz_member -> %w", err)
	}

	key := fmt.Sprintf(bizMembersCacheKeyFormat, bizCachePool, args.BizMember.BizId, args.BizMember.BizTenantId)

	err = its.refreshRpc.RefreshBizMembersCache(key, args.BizMember, rpcmodel.ActionInsert)
	if err != nil {
		log.Error("Error refresh biz_members cache: %v", zap.Error(err))
	}

	return nil
}

func (its *BizMemberService) RemoveBizMember(ctx context.Context, args *rpcmodel.BizMemberArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.DeleteBizMember(args.BizMember.BizId, args.BizMember.BizTenantId, args.BizMember.MemberId)
	if err != nil {
		return fmt.Errorf("remove biz_member -> %w", err)
	}

	key := fmt.Sprintf(bizMembersCacheKeyFormat, bizCachePool, args.BizMember.BizId, args.BizMember.BizTenantId)

	err = its.refreshRpc.RefreshBizMembersCache(key, args.BizMember, rpcmodel.ActionDelete)
	if err != nil {
		log.Error("Error refresh biz_members cache: %v", zap.Error(err))
	}

	return nil
}

func (its *BizMemberService) GetBizMembers(ctx context.Context, args *rpcmodel.BizMemberArgs, reply *rpcmodel.BizMembersReply) error {
	key := fmt.Sprintf(bizMembersCacheKeyFormat, bizMemberCachePool, args.BizMember.BizId, args.BizMember.BizTenantId)

	if members, exist := its.storageCache.Get(key); exist {
		reply.Members = members
		return nil
	}

	result, err, _ := singleGroup.Do(key, func() (interface{}, error) {
		members, err := its.database.GetBizMembers(args.BizMember.BizId, args.BizMember.BizTenantId)
		if err != nil {
			return nil, fmt.Errorf("get biz_members -> %w", err)
		}
		return members, nil
	})
	if err != nil {
		return fmt.Errorf("single group do -> %w", err)
	}

	members := result.([]*model.BizMember)
	reply.Members = members

	its.storageCache.Set(key, members)

	return nil
}
