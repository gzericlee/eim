package rpc

import (
	"context"
	"fmt"
	"time"

	"eim/internal/model"
	"eim/internal/redis"
	"eim/pkg/cache"
	"eim/pkg/cache/notify"
	"eim/util/log"
)

type BizMemberArgs struct {
	BizMember *model.BizMember
}

type BizMembersReply struct {
	Members []string
}

type BizMember struct {
	storageCache *cache.Cache
	redisManager *redis.Manager
}

func (its *BizMember) AddBizMember(ctx context.Context, args *BizMemberArgs, reply *EmptyReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.redisManager.AddBizMember(args.BizMember)
	if err != nil {
		return fmt.Errorf("save user -> %w", err)
	}

	key := fmt.Sprintf(cacheKeyFormat, bizCachePool, args.BizMember.BizId, args.BizMember.TenantId)
	err = notify.Del(bizMemberCachePool, key)
	if err != nil {
		return fmt.Errorf("del biz_members(%s) cache -> %w", key, err)
	}

	return nil
}

func (its *BizMember) GetBizMembers(ctx context.Context, args *BizMemberArgs, reply *BizMembersReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	key := fmt.Sprintf(cacheKeyFormat, bizMemberCachePool, args.BizMember.BizId, args.BizMember.TenantId)

	if cacheItem, exist := its.storageCache.Get(key); exist {
		reply.Members = cacheItem.([]string)
		return nil
	}

	result, err, _ := group.Do(key, func() (interface{}, error) {
		members, err := its.redisManager.GetBizMembers(args.BizMember.BizId, args.BizMember.TenantId)
		if err != nil {
			return nil, fmt.Errorf("get biz_members -> %w", err)
		}
		return members, nil
	})
	if err != nil {
		return fmt.Errorf("group do -> %w", err)
	}

	members := result.([]string)
	its.storageCache.Put(key, members)

	reply.Members = members

	return nil
}
