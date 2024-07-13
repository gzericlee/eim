package rpc

import (
	"context"
	"fmt"
	"time"

	rpcxetcdclient "github.com/rpcxio/rpcx-etcd/client"
	rpcxclient "github.com/smallnest/rpcx/client"

	"eim/internal/model"
)

type Client struct {
	devicePool    *rpcxclient.XClientPool
	messagePool   *rpcxclient.XClientPool
	bizPool       *rpcxclient.XClientPool
	tenantPool    *rpcxclient.XClientPool
	bizMemberPool *rpcxclient.XClientPool
	gatewayPool   *rpcxclient.XClientPool
	segmentPool   *rpcxclient.XClientPool
	refresherPool *rpcxclient.XClientPool
}

type EmptyArgs struct {
}

type EmptyReply struct {
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d1, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, deviceServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for device service -> %w", err)
	}

	d2, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, messageServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for message service -> %w", err)
	}

	d3, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, bizServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for biz service -> %w", err)
	}

	d4, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, bizMemberServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for biz_member service -> %w", err)
	}

	d5, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, gatewayServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for gateway service -> %w", err)
	}

	d6, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, segmentServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for segment service -> %w", err)
	}

	d7, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, refresherServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for refresher service -> %w", err)
	}

	d8, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, tenantServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for refresher service -> %w", err)
	}

	return &Client{
		devicePool:    rpcxclient.NewXClientPool(100, deviceServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d1, rpcxclient.DefaultOption),
		messagePool:   rpcxclient.NewXClientPool(100, messageServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d2, rpcxclient.DefaultOption),
		bizPool:       rpcxclient.NewXClientPool(100, bizServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d3, rpcxclient.DefaultOption),
		bizMemberPool: rpcxclient.NewXClientPool(100, bizMemberServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d4, rpcxclient.DefaultOption),
		gatewayPool:   rpcxclient.NewXClientPool(100, gatewayServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d5, rpcxclient.DefaultOption),
		segmentPool:   rpcxclient.NewXClientPool(100, segmentServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d6, rpcxclient.DefaultOption),
		refresherPool: rpcxclient.NewXClientPool(100, refresherServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d7, rpcxclient.DefaultOption),
		tenantPool:    rpcxclient.NewXClientPool(100, tenantServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, d8, rpcxclient.DefaultOption),
	}, nil
}

func (its *Client) InsertDevice(device *model.Device) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.devicePool.Get().Call(ctx, "InsertDevice", &DeviceArgs{Device: device}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertDevice -> %w", err)
	}
	return nil
}

func (its *Client) UpdateDevice(device *model.Device) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.devicePool.Get().Call(ctx, "UpdateDevice", &DeviceArgs{Device: device}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call UpdateDevice -> %w", err)
	}
	return nil
}

func (its *Client) GetDevice(userId, tenantId, deviceId string) (*model.Device, error) {
	reply := &DeviceReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.devicePool.Get().Call(ctx, "GetDevice", &UserArgs{UserId: userId, TenantId: tenantId, DeviceId: deviceId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetDevice -> %w", err)
	}
	return reply.Device, nil
}

func (its *Client) GetDevices(userId, tenantId string) ([]*model.Device, error) {
	reply := &DevicesReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.devicePool.Get().Call(ctx, "GetDevices", &UserArgs{UserId: userId, TenantId: tenantId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetDevices -> %w", err)
	}
	return reply.Devices, nil
}

func (its *Client) InsertMessage(message *model.Message) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.messagePool.Get().Call(ctx, "InsertMessage", &MessageArgs{Message: message}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertMessages -> %w", err)
	}
	return nil
}

func (its *Client) InsertMessages(messages []*model.Message) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.messagePool.Get().Call(ctx, "InsertMessages", &MessagesArgs{Messages: messages}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertMessages -> %w", err)
	}
	return nil
}

func (its *Client) SaveOfflineMessages(msgs []*model.Message, userId, deviceId string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.messagePool.Get().Call(ctx, "SaveOfflineMessages", &MessagesArgs{Messages: msgs, UserId: userId, DeviceId: deviceId}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call SaveOfflineMessages -> %w", err)
	}
	return nil
}

func (its *Client) RemoveOfflineMessages(msgIds []string, userId, deviceId string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.messagePool.Get().Call(ctx, "RemoveOfflineMessages", &MessageIdsArgs{MessageIds: msgIds, UserId: userId, DeviceId: deviceId}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RemoveOfflineMessages -> %w", err)
	}
	return nil
}

func (its *Client) GetOfflineMessagesCount(userId, deviceId string) (int64, error) {
	reply := &OfflineMessagesCountReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.messagePool.Get().Call(ctx, "GetOfflineMessagesCount", &OfflineMessagesArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return 0, fmt.Errorf("call GetOfflineMessagesCount -> %w", err)
	}
	return reply.Count, nil
}

func (its *Client) GetOfflineMessages(userId, deviceId string) ([]*model.Message, error) {
	reply := &MessagesReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.messagePool.Get().Call(ctx, "GetOfflineMessages", &OfflineMessagesArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetOfflineMessages -> %w", err)
	}
	return reply.Messages, nil
}

func (its *Client) InsertBiz(biz *model.Biz) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.bizPool.Get().Call(ctx, "InsertBiz", &BizArgs{Biz: biz}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertBiz -> %w", err)
	}
	return nil
}

func (its *Client) UpdateBiz(biz *model.Biz) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.bizPool.Get().Call(ctx, "UpdateBiz", &BizArgs{Biz: biz}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call UpdateBiz -> %w", err)
	}
	return nil
}

func (its *Client) GetBiz(bizId, tenantId string) (*model.Biz, error) {
	reply := &BizReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.bizPool.Get().Call(ctx, "GetBiz", &BizArgs{Biz: &model.Biz{BizId: bizId, TenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetBiz -> %w", err)
	}
	return reply.Biz, nil
}

func (its *Client) InsertTenant(tenant *model.Tenant) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.tenantPool.Get().Call(ctx, "InsertTenant", &TenantArgs{Tenant: tenant}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertTenant -> %w", err)
	}
	return nil
}

func (its *Client) UpdateTenant(tenant *model.Tenant) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.tenantPool.Get().Call(ctx, "UpdateTenant", &TenantArgs{Tenant: tenant}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call UpdateTenant -> %w", err)
	}
	return nil
}

func (its *Client) GetTenant(tenantId string) (*model.Tenant, error) {
	reply := &TenantReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.tenantPool.Get().Call(ctx, "GetTenant", &TenantArgs{Tenant: &model.Tenant{TenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetTenant -> %w", err)
	}
	return reply.Tenant, nil
}

func (its *Client) AddBizMember(member *model.BizMember) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.bizMemberPool.Get().Call(ctx, "AddBizMember", &BizMemberArgs{BizMember: member}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call AddBizMember -> %w", err)
	}
	return nil
}

func (its *Client) RemoveBizMember(member *model.BizMember) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.bizMemberPool.Get().Call(ctx, "RemoveBizMember", &BizMemberArgs{BizMember: member}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RemoveBizMember -> %w", err)
	}
	return nil
}

func (its *Client) GetBizMembers(bizId, tenantId string) ([]string, error) {
	reply := &BizMembersReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.bizMemberPool.Get().Call(ctx, "GetBizMembers", &BizMemberArgs{BizMember: &model.BizMember{BizId: bizId, BizTenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetBizMembers -> %w", err)
	}
	return reply.Members, nil
}

func (its *Client) RegisterGateway(gateway *model.Gateway, expiration time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.gatewayPool.Get().Call(ctx, "RegisterGateway", &GatewayArgs{Gateway: gateway, Expiration: expiration}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RegisterGateway -> %w", err)
	}
	return nil
}

func (its *Client) GetGateways() ([]*model.Gateway, error) {
	reply := &GatewaysReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.gatewayPool.Get().Call(ctx, "GetGateways", &EmptyArgs{}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetGateways -> %w", err)
	}
	return reply.Gateways, nil
}

func (its *Client) GetSegment(bizId, tenantId string) (*model.Segment, error) {
	reply := &SegmentReply{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.segmentPool.Get().Call(ctx, "GetSegment", &SegmentArgs{BizId: bizId, TenantId: tenantId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetSegment -> %w", err)
	}
	return reply.Segment, nil
}

func (its *Client) RefreshDevicesCache(key string, device *model.Device, action string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.refresherPool.Get().Broadcast(ctx, "RefreshDevicesCache", &RefreshDevicesArgs{Key: key, Device: device, Action: action}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RefreshDevicesCache -> %w", err)
	}
	return nil
}

func (its *Client) RefreshBizCache(key string, biz *model.Biz, action string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.refresherPool.Get().Broadcast(ctx, "RefreshBizCache", &RefreshBizArgs{Key: key, Biz: biz, Action: action}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RefreshBizsCache -> %w", err)
	}
	return nil
}

func (its *Client) RefreshBizMembersCache(key string, bizMember *model.BizMember, action string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.refresherPool.Get().Broadcast(ctx, "RefreshBizMembersCache", &RefreshBizMembersArgs{Key: key, BizMember: bizMember, Action: action}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RefreshBizMembersCache -> %w", err)
	}
	return nil
}

func (its *Client) RefreshTenantCache(key string, tenant *model.Tenant, action string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := its.refresherPool.Get().Broadcast(ctx, "RefreshTenantCache", &RefreshTenantArgs{Key: key, Tenant: tenant, Action: action}, &EmptyReply{})
	if err != nil {
		return fmt.Errorf("call RefreshTenantCache -> %w", err)
	}
	return nil
}
