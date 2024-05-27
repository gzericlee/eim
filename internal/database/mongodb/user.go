package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"eim/internal/model"
)

func (its *Repository) SaveUser(user *model.User) error {
	_, err := its.db.Collection("user").ReplaceOne(context.TODO(), bson.M{"user_id": user.UserId, "tenant_id": user.TenantId}, user, &options.ReplaceOptions{Upsert: &isTrue})
	if err != nil {
		return fmt.Errorf("upsert user -> %w", err)
	}
	return nil
}

func (its *Repository) GetUser(loginId, tenantId string) (*model.User, error) {
	var user *model.User
	err := its.db.Collection("user").FindOne(context.TODO(), bson.M{"login_id": loginId, "tenant_id": tenantId}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("find one user -> %w", err)
	}
	return user, nil
}
