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

type RefreshDevicesArgs struct {
	Key    string
	Device *model.Device
	Action string
}

type RefreshBizArgs struct {
	Key    string
	Biz    *model.Biz
	Action string
}

type RefreshBizMembersArgs struct {
	Key       string
	BizMember *model.BizMember
	Action    string
}

type Refresher struct {
	lock            *lock.KeyLock
	devicesCache    *cache.Cache[string, []*model.Device]
	bizCache        *cache.Cache[string, *model.Biz]
	bizMembersCache *cache.Cache[string, []string]
}

func (its *Refresher) RefreshDevicesCache(ctx context.Context, args *RefreshDevicesArgs, reply *EmptyReply) error {
	_, unlock := its.lock.Lock(args.Key, nil)
	defer unlock()

	switch args.Action {
	case ActionSave:
		{
			if devices, exist := its.devicesCache.Get(args.Key); exist {
				for i := range devices {
					if devices[i].DeviceId == args.Device.DeviceId {
						devices[i] = args.Device
						break
					}
				}
				its.devicesCache.Set(args.Key, devices)
			}
		}
	case ActionDelete:
		{
			if devices, exist := its.devicesCache.Get(args.Key); exist {
				for i := range devices {
					if devices[i].DeviceId == args.Device.DeviceId {
						devices = append(devices[:i], devices[i+1:]...)
						break
					}
				}
				its.devicesCache.Set(args.Key, devices)
			}
		}
	}

	return nil
}

func (its *Refresher) RefreshBizCache(ctx context.Context, args *RefreshBizArgs, reply *EmptyReply) error {
	_, unlock := its.lock.Lock(args.Key, nil)
	defer unlock()

	switch args.Action {
	case ActionSave:
		{
			its.bizCache.Set(args.Key, args.Biz)
		}
	}

	return nil
}

func (its *Refresher) RefreshBizMembersCache(ctx context.Context, args *RefreshBizMembersArgs, reply *EmptyReply) error {
	_, unlock := its.lock.Lock(args.Key, nil)
	defer unlock()

	switch args.Action {
	case ActionAdd:
		{
			if members, exist := its.bizMembersCache.Get(args.Key); exist {
				members = append(members, args.BizMember.MemberId)
				its.bizMembersCache.Set(args.Key, members)
			}
		}
	case ActionDelete:
		{
			if members, exist := its.bizMembersCache.Get(args.Key); exist {
				for i := range members {
					if members[i] == args.BizMember.MemberId {
						members = append(members[:i], members[i+1:]...)
						break
					}
				}
				its.bizMembersCache.Set(args.Key, members)
			}
		}
	}

	return nil
}
