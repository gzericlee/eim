package redis

import (
	"context"
	"fmt"

	"eim/internal/model"
)

func (its *Manager) SaveBizMember(member *model.BizMember) error {
	key := fmt.Sprintf("%s@%s:members", member.BizType, member.BizId)
	return its.redisClient.SAdd(context.Background(), key, member.UserId).Err()
}

func (its *Manager) GetBizMembers(bizType, bizId string) ([]string, error) {
	key := fmt.Sprintf("%s@%s:members", bizType, bizId)
	return its.redisClient.SMembers(context.Background(), key).Result()
}
