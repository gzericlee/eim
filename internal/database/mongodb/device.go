package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
)

func (its *Repository) SaveDevice(device *model.Device) error {
	_, err := its.db.Collection("device").ReplaceOne(context.TODO(), bson.M{"device_id": device.DeviceId}, device, &options.ReplaceOptions{Upsert: &isTrue})
	if err != nil {
		return fmt.Errorf("upsert device -> %w", err)
	}
	return nil
}

func (its *Repository) GetDevice(userId, deviceId string) (*model.Device, error) {
	var device *model.Device
	err := its.db.Collection("device").FindOne(context.TODO(), bson.M{"user_id": userId, "device_id": deviceId}).Decode(&device)
	if err != nil {
		return nil, fmt.Errorf("find one device -> %w", err)
	}
	return device, nil
}

func (its *Repository) GetDevices(userId string) ([]*model.Device, error) {
	var devices []*model.Device
	result, err := its.db.Collection("device").Find(context.TODO(), bson.M{"user_id": userId})
	if err != nil {
		return nil, fmt.Errorf("find user devices -> %w", err)
	}
	err = result.All(context.Background(), &devices)
	if err != nil {
		return nil, fmt.Errorf("find user devices -> %w", err)
	}
	return devices, nil
}
