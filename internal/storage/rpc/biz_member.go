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

type AppendBizMemberArgs struct {
	BizMember *model.BizMember
}

type GetBizMembersArgs struct {
	BizType string
	BizId   string
}

type GetBizMembersReply struct {
	Members []string
}

type BizMember struct {
	storageCache *cache.Cache
	redisManager *redis.Manager
}

func (its *BizMember) AppendBizMember(ctx context.Context, args *AppendBizMemberArgs, reply *EmptyReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.redisManager.AppendBizMember(args.BizMember)
	if err != nil {
		return fmt.Errorf("save user -> %w", err)
	}

	key := fmt.Sprintf("%s:%s:%s", userCachePool, args.BizMember.BizType, args.BizMember.BizId)
	err = notify.Del(bizMemberCachePool, key)
	if err != nil {
		return fmt.Errorf("del biz_members(%s) cache -> %w", key, err)
	}

	return nil
}

func (its *BizMember) GetBizMembers(ctx context.Context, args *GetBizMembersArgs, reply *GetBizMembersReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	key := fmt.Sprintf("%s:%s:%s", bizMemberCachePool, args.BizType, args.BizId)

	if cacheItem, exist := its.storageCache.Get(key); exist {
		reply.Members = cacheItem.([]string)
		return nil
	}

	result, err, _ := group.Do(key, func() (interface{}, error) {
		user, err := its.redisManager.GetBizMembers(args.BizType, args.BizId)
		if err != nil {
			return nil, fmt.Errorf("get biz_members -> %w", err)
		}
		return user, nil
	})
	if err != nil {
		return fmt.Errorf("group do -> %w", err)
	}

	members := result.([]string)
	its.storageCache.Put(key, members)

	reply.Members = members

	return nil
}
