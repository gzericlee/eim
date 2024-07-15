package pgsql

import (
	"fmt"

	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/internal/model/consts"
)

func (its *Repository) InsertTenant(tenant *model.Tenant) error {
	_, err := its.db.Insert(tenant)
	if err != nil {
		return fmt.Errorf("insert tenant -> %w", err)
	}
	return nil
}

func (its *Repository) UpdateTenant(tenant *model.Tenant) error {
	_, err := its.db.Where("tenant_id = ?", tenant.TenantId).Update(tenant)
	if err != nil {
		return fmt.Errorf("update tenant -> %w", err)
	}
	return nil
}

func (its *Repository) GetTenant(tenantId string) (*model.Tenant, error) {
	var tenant = &model.Tenant{}
	_, err := its.db.Where("tenant_id = ?", tenantId).Get(tenant)
	if err != nil {
		return nil, fmt.Errorf("select tenant -> %w", err)
	}
	return tenant, nil
}

func (its *Repository) EnableTenant(tenantId string) error {
	tenant, err := its.GetTenant(tenantId)
	if err != nil {
		return fmt.Errorf("get tenant -> %w", err)
	}
	tenant.State = consts.StatusEnabled
	err = its.UpdateTenant(tenant)
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
	err = its.UpdateTenant(tenant)
	if err != nil {
		return fmt.Errorf("disable tenant -> %w", err)
	}
	return nil
}
