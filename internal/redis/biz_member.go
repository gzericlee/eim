package redis

import (
	"context"
	"fmt"

	"eim/internal/model"
)

const (
	bizMemberKeyFormat = "biz:%s:%s:members"
)

func (its *Manager) AddBizMember(member *model.BizMember) error {
	key := fmt.Sprintf(bizMemberKeyFormat, member.BizId, member.TenantId)

	err := its.redisClient.SAdd(context.Background(), key, member.MemberId).Err()
	if err != nil {
		return fmt.Errorf("redis sadd(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) RemoveBizMember(member *model.BizMember) error {
	key := fmt.Sprintf(bizMemberKeyFormat, member.BizId, member.TenantId)

	err := its.redisClient.SRem(context.Background(), key, member.MemberId).Err()
	if err != nil {
		return fmt.Errorf("redis srem(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetBizMembers(bizId, tenantId string) ([]string, error) {
	key := fmt.Sprintf(bizMemberKeyFormat, bizId, tenantId)

	result, err := its.redisClient.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis smembers(%s) -> %w", key, err)
	}

	return result, nil
}
