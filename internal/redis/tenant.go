package redis

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	"eim/internal/model"
)

const (
	tenantKeyFormat = "tenant:%s"
)

func (its *Manager) SaveTenant(tenant *model.Tenant) error {
	key := fmt.Sprintf(tenantKeyFormat, tenant.TenantId)

	body, err := proto.Marshal(tenant)
	if err != nil {
		return fmt.Errorf("proto marshal -> %w", err)
	}

	err = its.redisClient.Set(context.Background(), key, body, 0).Err()
	if err != nil {
		return fmt.Errorf("redis set(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetTenant(tenantId string) (*model.Tenant, error) {
	key := fmt.Sprintf(tenantKeyFormat, tenantId)

	result, err := its.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get(%s) -> %w", key, err)
	}

	tenant := &model.Tenant{}
	err = proto.Unmarshal([]byte(result), tenant)
	if err != nil {
		return nil, fmt.Errorf("proto unmarshal -> %w", err)
	}

	return tenant, nil
}

func (its *Manager) EnableTenant(tenantId string) error {
	tenant, err := its.GetTenant(tenantId)
	if err != nil {
		return fmt.Errorf("get tenant -> %w", err)
	}
	tenant.State = model.Enabled
	err = its.SaveTenant(tenant)
	if err != nil {
		return fmt.Errorf("enable tenant -> %w", err)
	}

	return nil
}

func (its *Manager) DisableTenant(tenantId string) error {
	tenant, err := its.GetTenant(tenantId)
	if err != nil {
		return fmt.Errorf("get tenant -> %w", err)
	}
	tenant.State = model.Disabled
	err = its.SaveTenant(tenant)
	if err != nil {
		return fmt.Errorf("disable tenant -> %w", err)
	}

	return nil
}
