package redis

import (
	"context"
	"fmt"

	"eim/internal/model"
)

func (its *Manager) SaveObjectMember(member *model.ObjectMember) error {
	key := fmt.Sprintf("%s@%s:members", member.ObjectType, member.ObjectId)
	return its.redisClient.SAdd(context.Background(), key, member.UserId).Err()
}

func (its *Manager) GetObjectMembers(objectType, objectId string) ([]string, error) {
	key := fmt.Sprintf("%s@%s:members", objectType, objectId)
	return its.redisClient.SMembers(context.Background(), key).Result()
}
