package mongodb

import (
	"context"

	"eim/internal/model"
)

func (its *Repository) SaveMessage(message *model.Message) error {
	_, err := its.db.Collection("message").InsertOne(context.TODO(), message)
	return err
}
