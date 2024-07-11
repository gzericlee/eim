package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
	"eim/internal/model/consts"
)

func (its *Repository) SaveBiz(biz *model.Biz) error {
	_, err := its.db.Collection("biz").ReplaceOne(context.Background(), bson.M{"biz_id": biz.BizId, "tenant_id": biz.TenantId}, biz, &options.ReplaceOptions{Upsert: &isTrue})
	if err != nil {
		return fmt.Errorf("upsert biz -> %w", err)
	}
	return nil
}

func (its *Repository) GetBiz(bizId, tenantId string) (*model.Biz, error) {
	var biz *model.Biz
	err := its.db.Collection("biz").FindOne(context.Background(), bson.M{"biz_id": bizId, "tenant_id": tenantId}).Decode(&biz)
	if err != nil {
		return nil, fmt.Errorf("find one biz -> %w", err)
	}
	return biz, nil
}

func (its *Repository) EnableBiz(bizId, tenantId string) error {
	biz, err := its.GetBiz(bizId, tenantId)
	if err != nil {
		return fmt.Errorf("get biz -> %w", err)
	}
	biz.State = consts.StatusEnabled
	err = its.SaveBiz(biz)
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
	err = its.SaveBiz(biz)
	if err != nil {
		return fmt.Errorf("disable biz -> %w", err)
	}
	return nil
}

func (its *Repository) ListBizs(filter map[string]interface{}, limit, offset int64) ([]*model.Biz, int64, error) {
	total, err := its.db.Collection("biz").CountDocuments(context.Background(), bson.M(filter))
	if err != nil {
		return nil, 0, fmt.Errorf("count bizs -> %w", err)
	}

	var bizs []*model.Biz
	result, err := its.db.Collection("biz").Find(context.Background(), bson.M(filter), &options.FindOptions{Limit: &limit, Skip: &offset})
	if err != nil {
		return nil, total, fmt.Errorf("find bizs -> %w", err)
	}
	err = result.All(context.Background(), &bizs)
	if err != nil {
		return nil, total, fmt.Errorf("find bizs -> %w", err)
	}
	return bizs, total, nil
}
