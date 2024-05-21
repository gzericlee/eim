package redis

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	"eim/internal/model"
)

func (its *Manager) SaveUser(user *model.User) error {
	key := fmt.Sprintf("%s.%s.info", user.LoginId, user.TenantId)

	body, err := proto.Marshal(user)
	if err != nil {
		return fmt.Errorf("proto marshal -> %w", err)
	}

	err = its.redisClient.Set(context.Background(), key, body, 0).Err()
	if err != nil {
		return fmt.Errorf("redis set(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetUser(loginId, tenantId string) (*model.User, error) {
	key := fmt.Sprintf("%s.%s.info", loginId, tenantId)

	result, err := its.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get(%s) -> %w", key, err)
	}

	user := &model.User{}
	err = proto.Unmarshal([]byte(result), user)
	if err != nil {
		return nil, fmt.Errorf("proto unmarshal -> %w", err)
	}

	return user, nil
}
