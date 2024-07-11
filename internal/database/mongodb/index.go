package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (its *Repository) initIndexes() error {
	ctx := context.Background()

	messageIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"msg_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"seq_id": 1},
		},
	}

	segmentIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"biz_id": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	deviceIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"device_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"user_id": 1},
		},
	}

	tenantIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"tenant_id": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	bizIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"biz_id", 1}, {"tenant_id", 1}},
		},
	}

	bizMemberIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"biz_id", 1}, {"tenant_id", 1}, {"member_id", 1}},
		},
	}

	collections := map[string][]mongo.IndexModel{
		"message":    messageIndexes,
		"segment":    segmentIndexes,
		"device":     deviceIndexes,
		"tenant":     tenantIndexes,
		"biz":        bizIndexes,
		"biz_member": bizMemberIndexes,
	}

	for collectionName, indexes := range collections {
		for _, index := range indexes {
			_, _ = its.db.Collection(collectionName).Indexes().CreateOne(ctx, index)
		}
	}

	return nil
}
