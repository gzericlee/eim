package redis

import (
	"fmt"

	"eim/internal/model"
)

func (its *Manager) SaveGroupMember(member *model.GroupMember) error {
	key := fmt.Sprintf("%v:members", member.GroupId)
	return its.rdsClient.SAdd(key, member.UserId)
}

func (its *Manager) GetGroupMembers(groupId string) ([]string, error) {
	key := fmt.Sprintf("%v:members", groupId)
	return its.rdsClient.SMembers(key)
}
