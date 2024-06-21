package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
)

func (its *Repository) SaveBiz(biz *model.Biz) error {
	_, err := its.db.Collection("biz").ReplaceOne(context.TODO(), bson.M{"biz_id": biz.BizId, "tenant_id": biz.TenantId}, biz, &options.ReplaceOptions{Upsert: &isTrue})
	if err != nil {
		return fmt.Errorf("upsert biz -> %w", err)
	}
	return nil
}

func (its *Repository) GetBiz(bizId, tenantId string) (*model.Biz, error) {
	var biz *model.Biz
	err := its.db.Collection("biz").FindOne(context.TODO(), bson.M{"biz_id": bizId, "tenant_id": tenantId}).Decode(&biz)
	if err != nil {
		return nil, fmt.Errorf("find one biz -> %w", err)
	}
	return biz, nil
}
