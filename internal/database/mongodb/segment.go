package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"

	"eim/internal/model"
)

var (
	isTrue  = true
	isFalse = false
)

const (
	defaultStep = 1000
)

func (its *Repository) GetSegment(bizId string) (*model.Segment, error) {
	seg := &model.Segment{}
	err := its.db.Collection("segment").FindOne(context.Background(), bson.M{"biz_id": bizId}).Decode(seg)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		seg.BizId = bizId
		seg.MaxId = 0
		seg.Step = defaultStep
		seg.CreateAt = timestamppb.Now()
		seg.UpdateAt = timestamppb.Now()
	} else {
		seg.MaxId = seg.MaxId + int64(seg.Step)
		seg.UpdateAt = timestamppb.Now()
	}
	_, err = its.db.Collection("segment").ReplaceOne(context.Background(), bson.M{"biz_id": bizId}, seg, &options.ReplaceOptions{Upsert: &isTrue})
	return seg, err
}
