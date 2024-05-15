package mongodb

import (
	"context"
	"fmt"

	"eim/internal/model"
)

func (its *Repository) SaveMessage(message *model.Message) error {
	_, err := its.db.Collection("message").InsertOne(context.TODO(), message)
	if err != nil {
		return fmt.Errorf("insert message -> %w", err)
	}
	return nil
}
