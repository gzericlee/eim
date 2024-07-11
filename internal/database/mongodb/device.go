package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
)

func (its *Repository) SaveDevice(device *model.Device) error {
	_, err := its.db.Collection("device").ReplaceOne(context.Background(), bson.M{"device_id": device.DeviceId}, device, &options.ReplaceOptions{Upsert: &isTrue})
	if err != nil {
		return fmt.Errorf("upsert device -> %w", err)
	}
	return nil
}

func (its *Repository) GetDevice(userId, deviceId string) (*model.Device, error) {
	var device *model.Device
	err := its.db.Collection("device").FindOne(context.Background(), bson.M{"user_id": userId, "device_id": deviceId}).Decode(&device)
	if err != nil {
		return nil, fmt.Errorf("find one device -> %w", err)
	}
	return device, nil
}

func (its *Repository) GetDevicesByUser(userId, tenantId string) ([]*model.Device, error) {
	var devices []*model.Device
	result, err := its.db.Collection("device").Find(context.Background(), bson.M{"user_id": userId, "tenant_id": tenantId})
	if err != nil {
		return nil, fmt.Errorf("find user devices -> %w", err)
	}
	err = result.All(context.Background(), &devices)
	if err != nil {
		return nil, fmt.Errorf("find user devices -> %w", err)
	}
	return devices, nil
}

func (its *Repository) DeleteDevice(userId, tenantId, deviceId string) error {
	_, err := its.db.Collection("device").DeleteOne(context.Background(), bson.M{"user_id": userId, "tenant_id": tenantId, "device_id": deviceId})
	if err != nil {
		return fmt.Errorf("delete device -> %w", err)
	}
	return nil
}

func (its *Repository) ListDevices(filter map[string]interface{}, limit, offset int64) ([]*model.Device, int64, error) {
	total, err := its.db.Collection("device").CountDocuments(context.Background(), bson.M(filter))
	if err != nil {
		return nil, 0, fmt.Errorf("count bizs -> %w", err)
	}

	var devices []*model.Device
	result, err := its.db.Collection("device").Find(context.Background(), bson.M(filter), &options.FindOptions{Limit: &limit, Skip: &offset})
	if err != nil {
		return nil, total, fmt.Errorf("find devices -> %w", err)
	}
	err = result.All(context.Background(), &devices)
	if err != nil {
		return nil, total, fmt.Errorf("find devices -> %w", err)
	}
	return devices, total, nil
}