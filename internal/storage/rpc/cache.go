package rpc

import (
	"context"

	"eim/internal/model"
	"eim/pkg/cache"
	"eim/pkg/lock"
)

const (
	ActionSave   = "save"
	ActionDelete = "delete"
	ActionAdd    = "add"
)

type RefreshDeviceArgs struct {
	Key    string
	Device *model.Device
	Action string
}

type RefreshBizArgs struct {
	Key    string
	Biz    *model.Biz
	Action string
}

type RefreshBizMemberArgs struct {
	Key       string
	BizMember *model.BizMember
	Action    string
}

type Refresher struct {
	lock           *lock.KeyLock
	deviceCache    *cache.Cache
	bizCache       *cache.Cache
	bizMemberCache *cache.Cache
}

func (its *Refresher) RefreshDevicesCache(ctx context.Context, args *RefreshDeviceArgs, reply *EmptyReply) error {
	its.lock.Lock(args.Key)
	defer its.lock.Unlock(args.Key)

	switch args.Action {
	case ActionSave:
		{
			if cachedItem, exist := its.deviceCache.Get(args.Key); exist {
				devices := cachedItem.([]*model.Device)
				for i := range devices {
					if devices[i].DeviceId == args.Device.DeviceId {
						devices[i] = args.Device
						break
					}
				}
				its.deviceCache.Set(args.Key, devices)
			}
		}
	case ActionDelete:
		{
			if cachedItem, exist := its.deviceCache.Get(args.Key); exist {
				devices := cachedItem.([]*model.Device)
				for i := range devices {
					if devices[i].DeviceId == args.Device.DeviceId {
						devices = append(devices[:i], devices[i+1:]...)
						break
					}
				}
				its.deviceCache.Set(args.Key, devices)
			}
		}
	}

	return nil
}

func (its *Refresher) RefreshBizsCache(ctx context.Context, args *RefreshBizArgs, reply *EmptyReply) error {
	its.lock.Lock(args.Key)
	defer its.lock.Unlock(args.Key)

	switch args.Action {
	case ActionSave:
		{
			if cachedItem, exist := its.bizCache.Get(args.Key); exist {
				bizs := cachedItem.([]*model.Biz)
				for i := range bizs {
					if bizs[i].BizId == args.Biz.BizId && bizs[i].TenantId == args.Biz.TenantId {
						bizs[i] = args.Biz
						break
					}
				}
			}
		}
	case ActionDelete:
		{
			if cachedItem, exist := its.bizCache.Get(args.Key); exist {
				bizs := cachedItem.([]*model.Biz)
				for i := range bizs {
					if bizs[i].BizId == args.Biz.BizId && bizs[i].TenantId == args.Biz.TenantId {
						bizs = append(bizs[:i], bizs[i+1:]...)
						break
					}
				}
				its.bizCache.Set(args.Key, bizs)
			}
		}
	}

	return nil
}

func (its *Refresher) RefreshBizMembersCache(ctx context.Context, args *RefreshBizMemberArgs, reply *EmptyReply) error {
	its.lock.Lock(args.Key)
	defer its.lock.Unlock(args.Key)
	switch args.Action {
	case ActionAdd:
		{
			if cachedItem, exist := its.bizMemberCache.Get(args.Key); exist {
				members := cachedItem.([]string)
				members = append(members, args.BizMember.MemberId)
				its.bizMemberCache.Set(args.Key, members)
			}
		}
	case ActionDelete:
		{
			if cachedItem, exist := its.bizMemberCache.Get(args.Key); exist {
				members := cachedItem.([]string)
				for i := range members {
					if members[i] == args.BizMember.MemberId {
						members = append(members[:i], members[i+1:]...)
						break
					}
				}
				its.bizMemberCache.Set(args.Key, members)
			}
		}
	}

	return nil
}
