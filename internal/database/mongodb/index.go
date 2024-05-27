package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (its *Repository) initIndexes() error {
	ctx := context.TODO()

	_, err := its.db.Collection("device").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"device_id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create device(device_id) index -> %w", err)
	}

	_, err = its.db.Collection("device").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"user_id": 1},
		Options: options.Index().SetUnique(false),
	})
	if err != nil {
		return fmt.Errorf("create device(user_id) index -> %w", err)
	}

	_, err = its.db.Collection("device").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"device_id", 1}, {"user_id", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create device(device_id,user_id) index -> %w", err)
	}

	_, err = its.db.Collection("message").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"msg_id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create message(msg_id) index -> %w", err)
	}

	_, err = its.db.Collection("segment").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"biz_id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create segment(biz_id) index -> %w", err)
	}

	_, err = its.db.Collection("message").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"seq_id": 1},
	})
	if err != nil {
		return fmt.Errorf("create message(seq_id) index -> %w", err)
	}

	return nil
}
