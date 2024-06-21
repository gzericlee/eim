package rpc

import (
	"context"
	"fmt"
	"runtime"
	"time"

	rpcxetcdclient "github.com/rpcxio/rpcx-etcd/client"
	rpcxclient "github.com/smallnest/rpcx/client"

	"eim/internal/model"
)

type Client struct {
	devicePool    *rpcxclient.XClientPool
	messagePool   *rpcxclient.XClientPool
	userPool      *rpcxclient.XClientPool
	bizMemberPool *rpcxclient.XClientPool
	gatewayPool   *rpcxclient.XClientPool
	segmentPool   *rpcxclient.XClientPool
}

type EmptyArgs struct {
}

type EmptyReply struct {
}

func NewClient(etcdEndpoints []string) (*Client, error) {
	d1, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, deviceServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for device -> %w", err)
	}

	d2, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, messageServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for message -> %w", err)
	}

	d3, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, userServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for message -> %w", err)
	}

	d4, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, bizMemberServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for Message -> %w", err)
	}

	d5, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, gatewayServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for Message -> %w", err)
	}

	d6, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, segmentServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for message -> %w", err)
	}

	return &Client{
		devicePool:    rpcxclient.NewXClientPool(runtime.NumCPU()*10, deviceServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d1, rpcxclient.DefaultOption),
		messagePool:   rpcxclient.NewXClientPool(runtime.NumCPU()*10, messageServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d2, rpcxclient.DefaultOption),
		userPool:      rpcxclient.NewXClientPool(runtime.NumCPU()*10, userServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d3, rpcxclient.DefaultOption),
		bizMemberPool: rpcxclient.NewXClientPool(runtime.NumCPU()*10, bizMemberServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d4, rpcxclient.DefaultOption),
		gatewayPool:   rpcxclient.NewXClientPool(runtime.NumCPU()*10, gatewayServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d5, rpcxclient.DefaultOption),
		segmentPool:   rpcxclient.NewXClientPool(runtime.NumCPU()*10, segmentServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d6, rpcxclient.DefaultOption),
	}, nil
}

func (its *Client) SaveDevice(device *model.Device) error {
	err := its.devicePool.Get().Call(context.Background(), "SaveDevice", &DeviceArgs{Device: device}, nil)
	if err != nil {
		return fmt.Errorf("call SaveDevice -> %w", err)
	}
	return nil
}

func (its *Client) GetDevice(userId, deviceId string) (*model.Device, error) {
	reply := &DeviceReply{}
	err := its.devicePool.Get().Call(context.Background(), "GetDevice", &DeviceArgs{Device: &model.Device{UserId: userId, DeviceId: deviceId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetDevice -> %w", err)
	}
	return reply.Device, nil
}

func (its *Client) GetDevices(userId string) ([]*model.Device, error) {
	reply := &DevicesReply{}
	err := its.devicePool.Get().Call(context.Background(), "GetDevices", &DeviceArgs{Device: &model.Device{UserId: userId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetDevices -> %w", err)
	}
	return reply.Devices, nil
}

func (its *Client) SaveMessages(messages []*model.Message) error {
	err := its.messagePool.Get().Call(context.Background(), "SaveMessages", &MessagesArgs{Messages: messages}, nil)
	if err != nil {
		return fmt.Errorf("call SaveMessages -> %w", err)
	}
	return nil
}

func (its *Client) SaveOfflineMessages(msgs []*model.Message, userId, deviceId string) error {
	err := its.messagePool.Get().Call(context.Background(), "SaveOfflineMessages", &MessagesArgs{Messages: msgs, UserId: userId, DeviceId: deviceId}, nil)
	if err != nil {
		return fmt.Errorf("call SaveOfflineMessages -> %w", err)
	}
	return nil
}

func (its *Client) RemoveOfflineMessages(msgIds []string, userId, deviceId string) error {
	err := its.messagePool.Get().Call(context.Background(), "RemoveOfflineMessages", &MessageIdsArgs{MessageIds: msgIds, UserId: userId, DeviceId: deviceId}, nil)
	if err != nil {
		return fmt.Errorf("call RemoveOfflineMessages -> %w", err)
	}
	return nil
}

func (its *Client) GetOfflineMessagesCount(userId, deviceId string) (int64, error) {
	reply := &OfflineMessagesCountReply{}
	err := its.messagePool.Get().Call(context.Background(), "GetOfflineMessagesCount", &OfflineMessagesArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return 0, fmt.Errorf("call GetOfflineMessagesCount -> %w", err)
	}
	return reply.Count, nil
}

func (its *Client) GetOfflineMessages(userId, deviceId string) ([]*model.Message, error) {
	reply := &MessagesReply{}
	err := its.messagePool.Get().Call(context.Background(), "GetOfflineMessages", &OfflineMessagesArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetOfflineMessages -> %w", err)
	}
	return reply.Messages, nil
}

func (its *Client) SaveBiz(biz *model.Biz) error {
	err := its.userPool.Get().Call(context.Background(), "SaveBiz", &BizArgs{Biz: biz}, nil)
	if err != nil {
		return fmt.Errorf("call SaveBiz -> %w", err)
	}
	return nil
}

func (its *Client) GetBiz(bizId, tenantId string) (*model.Biz, error) {
	reply := &BizReply{}
	err := its.userPool.Get().Call(context.Background(), "GetBiz", &BizArgs{Biz: &model.Biz{BizId: bizId, TenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetBiz -> %w", err)
	}
	return reply.Biz, nil
}

func (its *Client) AddBizMember(member *model.BizMember) error {
	err := its.bizMemberPool.Get().Call(context.Background(), "AddBizMember", &BizMemberArgs{BizMember: member}, nil)
	if err != nil {
		return fmt.Errorf("call AddBizMember -> %w", err)
	}
	return nil
}

func (its *Client) GetBizMembers(bizId, tenantId string) ([]string, error) {
	reply := &BizMembersReply{}
	err := its.bizMemberPool.Get().Call(context.Background(), "GetBizMembers", &BizMemberArgs{BizMember: &model.BizMember{BizId: bizId, TenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetBizMembers -> %w", err)
	}
	return reply.Members, nil
}

func (its *Client) RegisterGateway(gateway *model.Gateway, expiration time.Duration) error {
	err := its.gatewayPool.Get().Call(context.Background(), "RegisterGateway", &GatewayArgs{Gateway: gateway, Expiration: expiration}, nil)
	if err != nil {
		return fmt.Errorf("call RegisterGateway -> %w", err)
	}
	return nil
}

func (its *Client) GetGateways() ([]*model.Gateway, error) {
	reply := &GatewaysReply{}
	err := its.bizMemberPool.Get().Call(context.Background(), "GetGateways", nil, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetGateways -> %w", err)
	}
	return reply.Gateways, nil
}

func (its *Client) GetSegment(bizId string) (*model.Segment, error) {
	reply := &SegmentReply{}
	err := its.segmentPool.Get().Call(context.Background(), "GetSegment", &SegmentArgs{BizId: bizId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetSegment -> %w", err)
	}
	return reply.Segment, nil
}
