package pgsql

import (
	"fmt"

	"eim/internal/model"
	"eim/internal/model/consts"
)

func (its *Repository) InsertBiz(biz *model.Biz) error {
	_, err := its.db.Insert(biz)
	if err != nil {
		return fmt.Errorf("insert biz -> %w", err)
	}
	return err
}

func (its *Repository) UpdateBiz(biz *model.Biz) error {
	_, err := its.db.Where("biz_id = ? AND tenant_id = ?", biz.BizId, biz.TenantId).Update(biz)
	if err != nil {
		return fmt.Errorf("update biz -> %w", err)
	}
	return err
}

func (its *Repository) GetBiz(bizId, tenantId string) (*model.Biz, error) {
	var biz = &model.Biz{}
	_, err := its.db.Where("biz_id = ? AND tenant_id = ?", bizId, tenantId).Get(biz)
	if err != nil {
		return nil, fmt.Errorf("select biz -> %w", err)
	}
	return biz, nil
}

func (its *Repository) EnableBiz(bizId, tenantId string) error {
	biz, err := its.GetBiz(bizId, tenantId)
	if err != nil {
		return fmt.Errorf("get biz -> %w", err)
	}
	biz.State = consts.StatusEnabled
	err = its.UpdateBiz(biz)
	if err != nil {
		return fmt.Errorf("enable biz -> %w", err)
	}
	return nil
}

func (its *Repository) DisableBiz(bizId, tenantId string) error {
	biz, err := its.GetBiz(bizId, tenantId)
	if err != nil {
		return fmt.Errorf("get biz -> %w", err)
	}
	biz.State = consts.StatusDisabled
	err = its.UpdateBiz(biz)
	if err != nil {
		return fmt.Errorf("disable biz -> %w", err)
	}
	return nil
}

func (its *Repository) ListBizs(filter map[string]interface{}, order []string, limit, offset int64) ([]*model.Biz, int64, error) {
	var bizs []*model.Biz

	query := its.db.Where("")
	for key, value := range filter {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	for _, by := range order {
		query = query.OrderBy(by)
	}

	total, err := query.Limit(int(limit), int(offset)).FindAndCount(&bizs)
	if err != nil {
		return nil, 0, fmt.Errorf("select bizs -> %w", err)
	}

	return bizs, total, nil
}
