package service

import (
	"context"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
	"github.com/gzericlee/eim/pkg/cache"
	"github.com/gzericlee/eim/pkg/lock"
)

type RefresherService struct {
	lock            *lock.KeyLock
	devicesCache    *cache.Cache[string, []*model.Device]
	bizCache        *cache.Cache[string, *model.Biz]
	bizMembersCache *cache.Cache[string, []string]
	tenantCache     *cache.Cache[string, *model.Tenant]
}

func NewRefresherService(lock *lock.KeyLock, devicesCache *cache.Cache[string, []*model.Device], bizCache *cache.Cache[string, *model.Biz], bizMembersCache *cache.Cache[string, []string], tenantCache *cache.Cache[string, *model.Tenant]) *RefresherService {
	return &RefresherService{
		lock:            lock,
		devicesCache:    devicesCache,
		bizCache:        bizCache,
		bizMembersCache: bizMembersCache,
		tenantCache:     tenantCache,
	}
}

func (its *RefresherService) RefreshDevicesCache(ctx context.Context, args *rpcmodel.RefreshDevicesArgs, reply *rpcmodel.EmptyReply) error {
	_, unlock := its.lock.Lock(args.Key, nil)
	defer unlock()

	switch args.Action {
	case rpcmodel.ActionInsert:
		{
			devices, _ := its.devicesCache.Get(args.Key)
			devices = append(devices, args.Device)
			its.devicesCache.Set(args.Key, devices)
		}
	case rpcmodel.ActionUpdate:
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
	case rpcmodel.ActionDelete:
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

func (its *RefresherService) RefreshBizCache(ctx context.Context, args *rpcmodel.RefreshBizArgs, reply *rpcmodel.EmptyReply) error {
	_, unlock := its.lock.Lock(args.Key, nil)
	defer unlock()

	switch args.Action {
	case rpcmodel.ActionInsert, rpcmodel.ActionUpdate:
		{
			its.bizCache.Set(args.Key, args.Biz)
		}
	}

	return nil
}

func (its *RefresherService) RefreshTenantCache(ctx context.Context, args *rpcmodel.RefreshTenantArgs, reply *rpcmodel.EmptyReply) error {
	_, unlock := its.lock.Lock(args.Key, nil)
	defer unlock()

	switch args.Action {
	case rpcmodel.ActionInsert, rpcmodel.ActionUpdate:
		{
			its.tenantCache.Set(args.Key, args.Tenant)
		}
	}

	return nil
}

func (its *RefresherService) RefreshBizMembersCache(ctx context.Context, args *rpcmodel.RefreshBizMembersArgs, reply *rpcmodel.EmptyReply) error {
	_, unlock := its.lock.Lock(args.Key, nil)
	defer unlock()

	switch args.Action {
	case rpcmodel.ActionInsert:
		{
			if members, exist := its.bizMembersCache.Get(args.Key); exist {
				members = append(members, args.BizMember.MemberId)
				its.bizMembersCache.Set(args.Key, members)
			}
		}
	case rpcmodel.ActionDelete:
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
