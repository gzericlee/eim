package redis

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	"eim/internal/model"
)

const (
	bizKeyFormat = "biz:%s:%s"
)

func (its *Manager) SaveBiz(biz *model.Biz) error {
	key := fmt.Sprintf(bizKeyFormat, biz.BizId, biz.TenantId)

	body, err := proto.Marshal(biz)
	if err != nil {
		return fmt.Errorf("proto marshal -> %w", err)
	}

	err = its.redisClient.Set(context.Background(), key, body, 0).Err()
	if err != nil {
		return fmt.Errorf("redis set(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetBiz(bizId, tenantId string) (*model.Biz, error) {
	key := fmt.Sprintf(bizKeyFormat, bizId, tenantId)

	result, err := its.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get(%s) -> %w", key, err)
	}

	biz := &model.Biz{}
	err = proto.Unmarshal([]byte(result), biz)
	if err != nil {
		return nil, fmt.Errorf("proto unmarshal -> %w", err)
	}

	return biz, nil
}
