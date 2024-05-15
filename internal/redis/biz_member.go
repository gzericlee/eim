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
		return fmt.Errorf("save biz member -> %w", err)
	}
	return nil
}

func (its *Manager) GetBizMembers(bizType, bizId string) ([]string, error) {
	key := fmt.Sprintf("%s.%s.members", bizType, bizId)
	result, err := its.redisClient.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("get biz members -> %w", err)
	}
	return result, nil
}
