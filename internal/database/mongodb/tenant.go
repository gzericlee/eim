package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
	"eim/internal/model/consts"
)

func (its *Repository) SaveTenant(tenant *model.Tenant) error {
	_, err := its.db.Collection("tenant").ReplaceOne(context.Background(), bson.M{"tenant_id": tenant.TenantId}, tenant, &options.ReplaceOptions{Upsert: &isTrue})
	if err != nil {
		return fmt.Errorf("upsert tenant -> %w", err)
	}
	return nil
}

func (its *Repository) GetTenant(tenantId string) (*model.Tenant, error) {
	var tenant *model.Tenant
	err := its.db.Collection("tenant").FindOne(context.Background(), bson.M{"tenant_id": tenantId}).Decode(&tenant)
	if err != nil {
		return nil, fmt.Errorf("find one tenant -> %w", err)
	}
	return tenant, nil
}

func (its *Repository) EnableTenant(tenantId string) error {
	tenant, err := its.GetTenant(tenantId)
	if err != nil {
		return fmt.Errorf("get tenant -> %w", err)
	}
	tenant.State = consts.StatusEnabled
	err = its.SaveTenant(tenant)
	if err != nil {
		return fmt.Errorf("enable tenant -> %w", err)
	}
	return nil
}

func (its *Repository) DisableTenant(tenantId string) error {
	tenant, err := its.GetTenant(tenantId)
	if err != nil {
		return fmt.Errorf("get tenant -> %w", err)
	}
	tenant.State = consts.StatusDisabled
	err = its.SaveTenant(tenant)
	if err != nil {
		return fmt.Errorf("disable tenant -> %w", err)
	}
	return nil
}
