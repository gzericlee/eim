package pgsql

import (
	"fmt"

	"eim/internal/model"
)

const (
	defaultStep = 1000
)

func (its *Repository) GetSegment(bizId, tenantId string) (*model.Segment, error) {
	var seg = &model.Segment{}
	has, err := its.db.Where("biz_id = ? AND tenant_id = ?", bizId, tenantId).Get(seg)
	if err != nil {
		return nil, fmt.Errorf("select segment -> %w", err)
	}
	if !has {
		seg = &model.Segment{
			BizId:    bizId,
			TenantId: tenantId,
			MaxId:    defaultStep,
			Step:     defaultStep,
		}
		_, err = its.db.Insert(seg)
		if err != nil {
			return nil, fmt.Errorf("insert segment -> %w", err)
		}
		return seg, nil
	}

	seg.MaxId += int64(seg.Step)

	_, err = its.db.Where("biz_id = ? AND tenant_id = ?", bizId, tenantId).Update(seg)
	if err != nil {
		return nil, fmt.Errorf("update segment -> %w", err)
	}

	return seg, nil
}
