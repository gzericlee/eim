package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gzericlee/eim/internal/model"
)

var (
	isTrue  = true
	isFalse = false
)

const (
	defaultStep = 1000
)

func (its *Repository) GetSegment(bizId, tenantId string) (*model.Segment, error) {
	var ctx = context.Background()
	var seg *model.Segment

	now := time.Now().Unix()

	err := its.db.Collection("segment").FindOne(ctx, bson.M{"biz_id": bizId, "tenant_id": tenantId}).Decode(&seg)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("find segment -> %w", err)
		}
		seg = &model.Segment{}
		seg.BizId = bizId
		seg.TenantId = tenantId
		seg.MaxId = 0
		seg.Step = defaultStep
		seg.CreatedAt = now
		seg.UpdatedAt = now
	} else {
		seg.MaxId = seg.MaxId + int64(seg.Step)
		seg.UpdatedAt = now
	}

	_, err = its.db.Collection("segment").ReplaceOne(ctx, bson.M{"biz_id": bizId, "tenant_id": tenantId}, seg, &options.ReplaceOptions{Upsert: &isTrue})
	if err != nil {
		return nil, fmt.Errorf("upsert segment -> %w", err)
	}

	return seg, nil
}
