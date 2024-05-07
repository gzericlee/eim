package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
)

func (its *Repository) SaveDevice(device *model.Device) error {
	_, err := its.db.Collection("device").ReplaceOne(context.TODO(), bson.M{"device_id": device.DeviceId}, device, &options.ReplaceOptions{Upsert: &isTrue})
	return err
}
