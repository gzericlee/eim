package mongodb

import (
	"context"

	"eim/internal/model"
)

func (its *Repository) SaveDevice(device *model.Device) error {
	_, err := its.db.Collection("device").InsertOne(context.TODO(), device)
	return err
}
