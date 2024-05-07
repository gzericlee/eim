package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"eim/internal/model"
)

func (its *Manager) RegisterGateway(gateway *model.Gateway, expiration time.Duration) error {
	body, err := proto.Marshal(gateway)
	if err != nil {
		return err
	}
	return its.redisClient.Set(context.Background(), fmt.Sprintf("gateway:%v", gateway.Ip), body, expiration).Err()
}

func (its *Manager) GetGateways() ([]*model.Gateway, error) {
	values, err := its.getAll(fmt.Sprintf("gateway:*"), 5000)
	if err != nil {
		return nil, err
	}
	var gateways []*model.Gateway
	for _, value := range values {
		gateway := &model.Gateway{}
		err = proto.Unmarshal([]byte(value), gateway)
		if err != nil {
			continue
		}
		gateways = append(gateways, gateway)
	}
	return gateways, err
}
