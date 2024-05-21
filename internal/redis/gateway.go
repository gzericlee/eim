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
	key := fmt.Sprintf("gateway.%s", gateway.Ip)

	body, err := proto.Marshal(gateway)
	if err != nil {
		return fmt.Errorf("proto marshal -> %w", err)
	}

	err = its.redisClient.Set(context.Background(), key, body, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis set(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetGateways() ([]*model.Gateway, error) {
	key := fmt.Sprintf("gateway.*")

	values, err := its.getAllValues(key)
	if err != nil {
		return nil, fmt.Errorf("redis getAllValues(%s) -> %w", key, err)
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
