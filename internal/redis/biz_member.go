package redis

import (
	"context"
	"fmt"

	"eim/internal/model"
)

func (its *Manager) SaveBizMember(member *model.BizMember) error {
	key := fmt.Sprintf("%s.%s.members", member.BizType, member.BizId)

	err := its.redisClient.SAdd(context.Background(), key, member.UserId).Err()
	if err != nil {
		return fmt.Errorf("redis sadd(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetBizMembers(bizType, bizId string) ([]string, error) {
	key := fmt.Sprintf("%s.%s.members", bizType, bizId)

	result, err := its.redisClient.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis smembers(%s) -> %w", key, err)
	}

	return result, nil
}
