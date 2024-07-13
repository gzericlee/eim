package pgsql

import (
	"fmt"

	"eim/internal/model"
)

func (its *Repository) InsertBizMember(member *model.BizMember) error {
	_, err := its.db.Insert(member)
	if err != nil {
		return fmt.Errorf("insert biz member -> %w", err)
	}
	return nil
}

func (its *Repository) DeleteBizMember(bizId, tenantId, memberId string) error {
	_, err := its.db.Where("biz_id = ? AND tenant_id = ? AND member_id = ?", bizId, tenantId, memberId).Delete()
	if err != nil {
		return fmt.Errorf("delete biz member -> %w", err)
	}
	return nil
}

func (its *Repository) GetBizMembers(bizId, tenantId string) ([]*model.BizMember, error) {
	var members []*model.BizMember
	err := its.db.Where("biz_id = ? AND biz_tenant_id = ?", bizId, tenantId).Find(&members)
	if err != nil {
		return nil, fmt.Errorf("select biz members -> %w", err)
	}
	return members, nil
}
