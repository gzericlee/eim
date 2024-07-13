package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
	"eim/internal/model/consts"
)

func (its *Repository) InsertBiz(biz *model.Biz) error {
	biz.CreatedAt = time.Now().Unix()
	_, err := its.db.Collection("biz").InsertOne(context.Background(), biz)
	if err != nil {
		return fmt.Errorf("insert biz -> %w", err)
	}
	return nil
}

func (its *Repository) UpdateBiz(biz *model.Biz) error {
	biz.UpdatedAt = time.Now().Unix()
	_, err := its.db.Collection("biz").UpdateOne(context.Background(), bson.M{"biz_id": biz.BizId, "tenant_id": biz.TenantId}, biz)
	if err != nil {
		return fmt.Errorf("update biz -> %w", err)
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
	biz.UpdatedAt = time.Now().Unix()
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
	biz.UpdatedAt = time.Now().Unix()
	err = its.UpdateBiz(biz)
	if err != nil {
		return fmt.Errorf("disable biz -> %w", err)
	}
	return nil
}

func (its *Repository) ListBizs(filter map[string]interface{}, order []string, limit, offset int64) ([]*model.Biz, int64, error) {
	total, err := its.db.Collection("biz").CountDocuments(context.Background(), bson.M(filter))
	if err != nil {
		return nil, 0, fmt.Errorf("count bizs -> %w", err)
	}

	var orderBy = map[string]interface{}{}
	for _, by := range order {
		col := strings.Split(by, " ")[0]
		orderBy[col] = -1
		sort := strings.Split(by, " ")[1]
		if strings.EqualFold(sort, "asc") {
			orderBy[col] = 1
		}
	}

	var bizs []*model.Biz
	result, err := its.db.Collection("biz").Find(context.Background(), bson.M(filter), &options.FindOptions{Limit: &limit, Skip: &offset, Sort: bson.M(orderBy)})
	if err != nil {
		return nil, total, fmt.Errorf("find bizs -> %w", err)
	}
	err = result.All(context.Background(), &bizs)
	if err != nil {
		return nil, total, fmt.Errorf("find bizs -> %w", err)
	}

	return bizs, total, nil
}
