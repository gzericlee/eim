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
		devicePool:    rpcxclient.NewXClientPool(runtime.NumCPU()*2, deviceServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d1, rpcxclient.DefaultOption),
		messagePool:   rpcxclient.NewXClientPool(runtime.NumCPU()*2, messageServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d2, rpcxclient.DefaultOption),
		userPool:      rpcxclient.NewXClientPool(runtime.NumCPU()*2, userServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d3, rpcxclient.DefaultOption),
		bizMemberPool: rpcxclient.NewXClientPool(runtime.NumCPU()*2, bizMemberServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d4, rpcxclient.DefaultOption),
		gatewayPool:   rpcxclient.NewXClientPool(runtime.NumCPU()*2, gatewayServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d5, rpcxclient.DefaultOption),
		segmentPool:   rpcxclient.NewXClientPool(runtime.NumCPU()*2, segmentServicePath, rpcxclient.Failover, rpcxclient.ConsistentHash, d6, rpcxclient.DefaultOption),
	}, nil
}

func (its *Client) SaveDevice(device *model.Device) error {
	err := its.devicePool.Get().Call(context.Background(), "SaveDevice", &SaveDeviceArgs{Device: device}, nil)
	if err != nil {
		return fmt.Errorf("call save device -> %w", err)
	}
	return nil
}

func (its *Client) GetDevice(userId, deviceId string) (*model.Device, error) {
	reply := &DeviceReply{}
	err := its.devicePool.Get().Call(context.Background(), "GetDevice", &GetDeviceArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call get device -> %w", err)
	}
	return reply.Device, nil
}

func (its *Client) GetDevices(userId string) ([]*model.Device, error) {
	reply := &DevicesReply{}
	err := its.devicePool.Get().Call(context.Background(), "GetDevices", &GetDevicesArgs{UserId: userId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call get devices -> %w", err)
	}
	return reply.Devices, nil
}

func (its *Client) SaveMessages(messages []*model.Message) error {
	err := its.messagePool.Get().Call(context.Background(), "SaveMessages", &MessageArgs{Messages: messages}, nil)
	if err != nil {
		return fmt.Errorf("call save messages -> %w", err)
	}
	return nil
}

func (its *Client) SaveOfflineMessageIds(msgIds []interface{}, userId, deviceId string) error {
	err := its.messagePool.Get().Call(context.Background(), "SaveOfflineMessageIds", &MessageIdsArgs{MsgIds: msgIds, UserId: userId, DeviceId: deviceId}, nil)
	if err != nil {
		return fmt.Errorf("call save offline message ids -> %w", err)
	}
	return nil
}

func (its *Client) RemoveOfflineMessageIds(msgIds []interface{}, userId, deviceId string) error {
	err := its.messagePool.Get().Call(context.Background(), "RemoveOfflineMessageIds", &MessageIdsArgs{MsgIds: msgIds, UserId: userId, DeviceId: deviceId}, nil)
	if err != nil {
		return fmt.Errorf("call remove offline message ids -> %w", err)
	}
	return nil
}

func (its *Client) GetOfflineMessagesCount(userId, deviceId string) (int64, error) {
	reply := &OfflineMessagesCountReply{}
	err := its.messagePool.Get().Call(context.Background(), "GetOfflineMessagesCount", &OfflineMessagesArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return 0, fmt.Errorf("call get offline messages count -> %w", err)
	}
	return reply.Count, nil
}

func (its *Client) GetOfflineMessages(userId, deviceId string) ([]*model.Message, error) {
	reply := &MessagesReply{}
	err := its.messagePool.Get().Call(context.Background(), "GetOfflineMessages", &OfflineMessagesArgs{UserId: userId, DeviceId: deviceId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call get offline messages -> %w", err)
	}
	return reply.Messages, nil
}

func (its *Client) SaveUser(user *model.User) error {
	err := its.userPool.Get().Call(context.Background(), "SaveUser", &UserArgs{User: user}, nil)
	if err != nil {
		return fmt.Errorf("call save user -> %w", err)
	}
	return nil
}

func (its *Client) GetUser(loginId, tenantId string) (*model.User, error) {
	reply := &UserReply{}
	err := its.userPool.Get().Call(context.Background(), "GetUser", &UserArgs{User: &model.User{LoginId: loginId, TenantId: tenantId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call get user -> %w", err)
	}
	return reply.User, nil
}

func (its *Client) AppendBizMember(member *model.BizMember) error {
	err := its.bizMemberPool.Get().Call(context.Background(), "AppendBizMember", &AppendBizMemberArgs{BizMember: member}, nil)
	if err != nil {
		return fmt.Errorf("call append biz_member -> %w", err)
	}
	return nil
}

func (its *Client) GetBizMembers(bizType, bizId string) ([]string, error) {
	reply := &GetBizMembersReply{}
	err := its.bizMemberPool.Get().Call(context.Background(), "GetBizMembers", &GetBizMembersArgs{BizType: bizType, BizId: bizId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call append biz_member -> %w", err)
	}
	return reply.Members, nil
}

func (its *Client) RegisterGateway(gateway *model.Gateway, expiration time.Duration) error {
	err := its.gatewayPool.Get().Call(context.Background(), "RegisterGateway", &RegisterGatewayArgs{Gateway: gateway, Expiration: expiration}, nil)
	if err != nil {
		return fmt.Errorf("call register gateway -> %w", err)
	}
	return nil
}

func (its *Client) GetGateways() ([]*model.Gateway, error) {
	reply := &GetGatewaysReply{}
	err := its.bizMemberPool.Get().Call(context.Background(), "GetBizMembers", nil, reply)
	if err != nil {
		return nil, fmt.Errorf("call get gateways -> %w", err)
	}
	return reply.Gateways, nil
}

func (its *Client) GetSegment(bizId string) (*model.Segment, error) {
	reply := &SegmentReply{}
	err := its.segmentPool.Get().Call(context.Background(), "GetSegment", &GetSegmentArgs{BizId: bizId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call get Segment -> %w", err)
	}
	return reply.Segment, nil
}
