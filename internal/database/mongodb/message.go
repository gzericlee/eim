package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"eim/internal/model"
)

func (its *Repository) SaveMessage(message *model.Message) error {
	_, err := its.db.Collection("message").InsertOne(context.TODO(), message)
	if err != nil {
		return fmt.Errorf("insert message -> %w", err)
	}
	return nil
}

func (its *Repository) SaveMessages(messages []*model.Message) error {
	var objs []interface{}
	for _, message := range messages {
		objs = append(objs, message)
	}

	_, err := its.db.Collection("message").InsertMany(context.TODO(), objs)
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
