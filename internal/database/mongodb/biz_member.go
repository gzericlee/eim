package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
)

func (its *Repository) InsertBizMember(member *model.BizMember) error {
	_, err := its.db.Collection("biz_member").InsertOne(context.Background(), member, &options.InsertOneOptions{})
	if err != nil {
		return fmt.Errorf("insert biz member -> %w", err)
	}
	return nil
}

func (its *Repository) DeleteBizMember(bizId, tenantId, memberId string) error {
	_, err := its.db.Collection("biz_member").DeleteOne(context.Background(), bson.M{"biz_id": bizId, "tenant_id": tenantId, "member_id": memberId})
	if err != nil {
		return fmt.Errorf("delete biz member -> %w", err)
	}
	return nil
}

func (its *Repository) GetBizMembers(bizId, tenantId string) ([]*model.BizMember, error) {
	var members []*model.BizMember

	cursor, err := its.db.Collection("biz_member").Find(context.Background(), bson.M{"biz_id": bizId, "tenant_id": tenantId})
	if err != nil {
		return nil, fmt.Errorf("find biz members -> %w", err)
	}

	err = cursor.All(context.Background(), &members)
	if err != nil {
		return nil, fmt.Errorf("cursor all -> %w", err)
	}

	return members, nil
}
