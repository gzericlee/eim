package mongodb

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
)

func (its *Repository) InsertMessage(message *model.Message) error {
	_, err := its.db.Collection("message").InsertOne(context.Background(), message)
	if err != nil {
		return fmt.Errorf("insert message -> %w", err)
	}
	return nil
}

func (its *Repository) InsertMessages(messages []*model.Message) error {
	var objs []interface{}
	for _, message := range messages {
		objs = append(objs, message)
	}

	_, err := its.db.Collection("message").InsertMany(context.Background(), objs)
	if err != nil {
		return fmt.Errorf("insert messages -> %w", err)
	}

	return nil
}

func (its *Repository) GetMessagesByIds(msgIds []int64) ([]*model.Message, error) {
	var ctx = context.Background()
	var messages []*model.Message

	cursor, err := its.db.Collection("message").Find(ctx, bson.M{"msg_id": bson.M{"$in": msgIds}})
	if err != nil {
		return nil, fmt.Errorf("get messages by ids -> %w", err)
	}

	if err = cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("cursor all -> %w", err)
	}

	return messages, nil
}

func (its *Repository) ListHistoryMessages(filter map[string]interface{}, order []string, minSeq, maxSeq, limit, offset int64) ([]*model.Message, error) {
	var messages []*model.Message

	bm := bson.M(filter)

	var orderBy = map[string]interface{}{}
	for _, by := range order {
		col := strings.Split(by, " ")[0]
		orderBy[col] = -1
		sort := strings.Split(by, " ")[1]
		if strings.EqualFold(sort, "asc") {
			orderBy[col] = 1
		}
	}

	if minSeq > 0 && maxSeq > 0 && maxSeq > minSeq {
		bm["$and"] = []interface{}{
			bson.M{"seq_id": bson.M{"$gte": minSeq}},
			bson.M{"seq_id": bson.M{"$lte": maxSeq}},
		}
	}

	result, err := its.db.Collection("message").Find(context.Background(), bm, &options.FindOptions{Limit: &limit, Skip: &offset, Sort: bson.M(orderBy)})
	if err != nil {
		return nil, fmt.Errorf("find messages -> %w", err)
	}

	err = result.All(context.Background(), &messages)
	if err != nil {
		return nil, fmt.Errorf("find messages -> %w", err)
	}

	return messages, nil
}
