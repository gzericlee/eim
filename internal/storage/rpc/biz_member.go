package rpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/redis"
	"eim/pkg/cache"
	"eim/pkg/log"
)

type BizMemberArgs struct {
	BizMember *model.BizMember
}

type BizMembersReply struct {
	Members []string
}

type BizMember struct {
	storageCache *cache.Cache[string, []string]
	redisManager *redis.Manager
}

func (its *BizMember) AddBizMember(ctx context.Context, args *BizMemberArgs, reply *EmptyReply) error {
	err := its.redisManager.AddBizMember(args.BizMember)
	if err != nil {
		return fmt.Errorf("add biz_member -> %w", err)
	}

	key := fmt.Sprintf(cacheKeyFormat, bizCachePool, args.BizMember.BizId, args.BizMember.BizTenantId)

	err = storageRpc.RefreshBizMembersCache(key, args.BizMember, ActionAdd)
	if err != nil {
		log.Error("Error refresh biz_members cache: %v", zap.Error(err))
	}

	return nil
}

func (its *BizMember) RemoveBizMember(ctx context.Context, args *BizMemberArgs, reply *EmptyReply) error {
	err := its.redisManager.RemoveBizMember(args.BizMember)
	if err != nil {
		return fmt.Errorf("remove biz_member -> %w", err)
	}

	key := fmt.Sprintf(cacheKeyFormat, bizCachePool, args.BizMember.BizId, args.BizMember.BizTenantId)

	err = storageRpc.RefreshBizMembersCache(key, args.BizMember, ActionDelete)
	if err != nil {
		log.Error("Error refresh biz_members cache: %v", zap.Error(err))
	}

	return nil
}

func (its *BizMember) GetBizMembers(ctx context.Context, args *BizMemberArgs, reply *BizMembersReply) error {
	key := fmt.Sprintf(cacheKeyFormat, bizMemberCachePool, args.BizMember.BizId, args.BizMember.BizTenantId)

	if members, exist := its.storageCache.Get(key); exist {
		reply.Members = members
		return nil
	}

	result, err, _ := singleGroup.Do(key, func() (interface{}, error) {
		members, err := its.redisManager.GetBizMembers(args.BizMember.BizId, args.BizMember.BizTenantId)
		if err != nil {
			return nil, fmt.Errorf("get biz_members -> %w", err)
		}
		return members, nil
	})
	if err != nil {
		return fmt.Errorf("biz_member single group do -> %w", err)
	}

	members := result.([]string)
	reply.Members = members

	its.storageCache.Set(key, members)

	return nil
}
