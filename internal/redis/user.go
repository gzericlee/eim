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
		return err
	}
	return its.redisClient.Set(context.Background(), key, body, 0).Err()
}

func (its *Manager) GetUser(loginId, tenantId string) (*model.User, error) {
	key := fmt.Sprintf("%s.%s.info", loginId, tenantId)
	result, err := its.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = proto.Unmarshal([]byte(result), user)
	return user, err
}
