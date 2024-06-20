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

type UserArgs struct {
	User *model.User
}

type UserReply struct {
	User *model.User
}

type User struct {
	storageCache *cache.Cache
	redisManager *redis.Manager
	database     database.IDatabase
}

func (its *User) SaveUser(ctx context.Context, args *UserArgs, reply *EmptyReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	err := its.redisManager.SaveUser(args.User)
	if err != nil {
		return fmt.Errorf("save user -> %w", err)
	}

	key := fmt.Sprintf("%s:%s:%s", userCachePool, args.User.UserId, args.User.TenantId)
	err = notify.Del(userCachePool, key)
	if err != nil {
		return fmt.Errorf("del user(%s) cache -> %w", key, err)
	}

	return nil
}

func (its *User) GetUser(ctx context.Context, args *UserArgs, reply *UserReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	key := fmt.Sprintf("%s:%s:%s", userCachePool, args.User.LoginId, args.User.TenantId)

	if cacheItem, exist := its.storageCache.Get(key); exist {
		reply.User = cacheItem.(*model.User)
		return nil
	}

	result, err, _ := group.Do(key, func() (interface{}, error) {
		user, err := its.redisManager.GetUser(args.User.LoginId, args.User.TenantId)
		if err != nil {
			return nil, fmt.Errorf("get user -> %w", err)
		}
		return user, nil
	})
	if err != nil {
		return fmt.Errorf("group do -> %w", err)
	}

	user := result.(*model.User)
	its.storageCache.Put(key, user)

	reply.User = user

	return nil
}
