package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/util/log"
)

func (its *Manager) RegisterGateway(gateway *model.Gateway, expiration time.Duration) error {
	body, err := proto.Marshal(gateway)
	if err != nil {
		return fmt.Errorf("proto marshal -> %w", err)
	}
	err = its.redisClient.Set(context.Background(), fmt.Sprintf("gateway.%v", gateway.Ip), body, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis set -> %w", err)
	}
	return nil
}

func (its *Manager) GetGateways() ([]*model.Gateway, error) {
	values, err := its.getAll(fmt.Sprintf("gateway.*"), 5000)
	if err != nil {
		return nil, fmt.Errorf("redis getAll -> %w", err)
	}
	var gateways []*model.Gateway
	for _, value := range values {
		gateway := &model.Gateway{}
		err = proto.Unmarshal([]byte(value), gateway)
		if err != nil {
			log.Error("Error proto unmarshal. Drop it", zap.Error(err))
			continue
		}
		gateways = append(gateways, gateway)
	}
	return gateways, nil
}
