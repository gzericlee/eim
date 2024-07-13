package pgsql

import (
	"fmt"

	"xorm.io/xorm"

	"eim/internal/model"
)

type Repository struct {
	db *xorm.EngineGroup
}

func NewRepository(db *xorm.EngineGroup) (*Repository, error) {
	err := db.Sync(new(model.Device), new(model.Message), new(model.Segment), new(model.Tenant), new(model.Biz), new(model.BizMember))
	if err != nil {
		return nil, fmt.Errorf("sync table -> %w", err)
	}

	repository := &Repository{db: db}
	return repository, nil
}
